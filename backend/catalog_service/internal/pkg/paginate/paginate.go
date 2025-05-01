//* Пакет paginate предоставляет реализацию курсорной пагинации с шифрованием
//* курсорных токенов через интерфейс TokenCoder (например, AES-GCM).
//*
//* Основные элементы:
//*   - TokenCoder: интерфейс для шифрования/дешифрования байтовых данных
//*   - Paginator[T]: универсальный пагинатор для структур с набором полей для сортировки
//*
//* Пример использования:
//*
//*   coder, _ := NewEncryptor(key, nil)
//*   paginator, _ := NewPaginator[MyStruct](conn, dialect, coder)
//*   paginator = paginator.WithTable("my_table").WithColumns("id", "created_at").WithLimit(50)
//*   items, nextToken, err := paginator.Paginate(ctx, prevToken)

package paginate

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// * init регистрирует необходимые типы в пакете encoding/gob,
// * чтобы значения интерфейсных типов (например, time.Time и любых других
// * зарегистрированных вами типов) могли корректно сериализоваться.
// *
// * Без такой регистрации gob не сможет закодировать интерфейсные значения,
// * конкретные типы которых ему не известны.
func init() {
	gob.Register(time.Time{})
	gob.Register(uuid.UUID{})
	gob.Register([]byte{})
}

// * TokenCoder задаёт интерфейс для шифрования и дешифрования курсорных токенов.
// * Реализации могут использовать AES-GCM, HMAC или другие методы.
type TokenCoder interface {
	//* Encrypt кодирует данные в строковый токен.
	Encrypt(data []byte) (string, error)
	//* Decrypt декодирует строковый токен обратно в байты данных.
	Decrypt(token string) ([]byte, error)
}

// * Paginator выполняет курсорную пагинацию над таблицей базы данных.
// * T — любая структура (struct) с экспортируемыми полями, помеченными тэгом `db` или совпадающими
// * по имени с колонками таблицы. Поля используются для сканирования и извлечения значений курсора.
type Paginator[T any] struct {
	conn    *pgxpool.Pool       //* соединение с базой данных
	dialect goqu.DialectWrapper //* диалект SQL (Postgres, MySQL и т.д.)
	coder   TokenCoder          //* шифратор/дешифратор токенов
	limit   uint                //* максимальный размер страницы
	cols    []string            //* колонки для сортировки и курсора
	table   string              //* имя таблицы
}

// * DefaultLimit задаёт значение лимита по умолчанию, если не указано иное.
const DefaultLimit = 20

// * NewPaginator создаёт Paginator[T] с соединением, диалектом и шифратором.
// * Если conn или coder равны nil, возвращается ошибка.
// * Лимит устанавливается в DefaultLimit.
func NewPaginator[T any](conn *pgxpool.Pool, dialect goqu.DialectWrapper, coder TokenCoder) (*Paginator[T], error) {
	if conn == nil {
		return nil, errors.New("db connection is nil")
	}
	if coder == nil {
		return nil, errors.New("coder is nil")
	}

	p := &Paginator[T]{
		conn:    conn,
		dialect: dialect,
		coder:   coder,
		limit:   DefaultLimit,
	}
	return p, nil
}

// * WithTable устанавливает имя таблицы для пагинации.
func (p *Paginator[T]) WithTable(table string) *Paginator[T] {
	p.table = table
	return p
}

// * WithLimit задаёт максимальное число записей на страницу (если >0).
func (p *Paginator[T]) WithLimit(limit uint) *Paginator[T] {
	if limit > 0 {
		p.limit = limit
	}
	return p
}

// * WithColumns задаёт колонки, по которым производится упорядочивание и курсор.
// * Порядок колонок важен: первым идёт основной ключ, далее дополнительные.
func (p *Paginator[T]) WithColumns(cols ...string) *Paginator[T] {
	//* создаём копию слайса cols
	p.cols = append([]string(nil), cols...)
	return p
}

// * Paginate выполняет запрос к БД и возвращает элементы, новый токен курсора и ошибку.
// * token: предыдущий токен; если пусто — возвращает первую страницу.
// * Контекст ctx используется для отмены и таймаутов.
func (p *Paginator[T]) Paginate(ctx context.Context, token string) ([]T, string, error) {
	if p.limit == 0 {
		return nil, "", errors.New("limit must be > 0")
	}
	if len(p.cols) == 0 {
		return nil, "", errors.New("no columns specified for ordering")
	}
	if p.table == "" {
		return nil, "", errors.New("table name not set")
	}

	//* Строим набор данных
	ds, err := p.buildDataset(token)
	if err != nil {
		return nil, "", err
	}
	sqlStr, args, err := ds.ToSQL()
	if err != nil {
		return nil, "", fmt.Errorf("SQL build: %w", err)
	}

	//* Выполняем запрос
	rows, err := p.conn.Query(ctx, sqlStr, args...)
	// TODO: проверить запрос
	if err != nil {
		return nil, "", fmt.Errorf("query execution: %w", err)
	}
	defer rows.Close()

	//* Сканируем строки в срез структур T
	items, err := scanRows[T](rows)
	if err != nil {
		return nil, "", err
	}
	if len(items) == 0 {
		return items, "", nil
	}

	// * Проверяем существует ли следующая страница
	// * Если её нет то возвращаем пустой токен
	hasNext := len(items) > int(p.limit)
	if hasNext {
		//* Извлекаем значения последнего элемента для курсора
		vals, err := extractValuesReflect(items[len(items)-1], p.cols)
		if err != nil {
			return nil, "", fmt.Errorf("extract cursor values: %w", err)
		}

		//* Кодируем курсор в бинарный вид
		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(vals); err != nil {
			return nil, "", fmt.Errorf("encoding cursor: %w", err)
		}
		//* Шифруем и кодируем токен
		tok, err := p.coder.Encrypt(buf.Bytes())
		if err != nil {
			return nil, "", fmt.Errorf("encrypting cursor: %w", err)
		}
		return items[:p.limit], tok, nil
	}

	return items, "", nil
}

// * buildDataset конструирует goqu.SelectDataset с учётом токена.
// * Если token непустой — добавляет условие WHERE для курсора.
func (p *Paginator[T]) buildDataset(token string) (*goqu.SelectDataset, error) {
	//* Формируем выражения ORDER BY col DESC
	exprs := make([]exp.OrderedExpression, len(p.cols))
	for i, c := range p.cols {
		exprs[i] = goqu.C(c).Desc()
	}

	//* SELECT * FROM table ORDER BY ... LIMIT ...
	ds := p.dialect.From(p.table).Order(exprs...).Limit(p.limit + 1)
	if token != "" {
		//* Дешифруем и декодируем курсор
		raw, err := p.coder.Decrypt(token)
		if err != nil {
			return nil, fmt.Errorf("token decrypt: %w", err)
		}
		var vals []interface{}
		if err := gob.NewDecoder(bytes.NewReader(raw)).Decode(&vals); err != nil {
			return nil, fmt.Errorf("decoding cursor: %w", err)
		}
		//* Строим условие WHERE c1 < v1 OR (c1 = v1 AND c2 < v2) ...
		expr, err := cursorConditionExpr(p.cols, vals)
		if err != nil {
			return nil, err
		}
		ds = ds.Where(expr)
	}
	return ds, nil
}

// * cursorConditionExpr строит выражение для курсора: последовательность OR/AND по колонкам.
func cursorConditionExpr(cols []string, vals []interface{}) (goqu.Expression, error) {
	if len(cols) != len(vals) {
		return nil, errors.New("cols and vals length mismatch")
	}
	var expr goqu.Expression
	for i := range cols {
		//* условие cols[i] < vals[i]
		lt := goqu.Ex{cols[i]: goqu.Op{"lt": vals[i]}}
		var curr goqu.Expression = lt
		if i > 0 {
			//* условие cols[0]=vals[0] AND ... cols[i-1]=vals[i-1]
			eq := goqu.Ex{}
			for j := 0; j < i; j++ {
				eq[cols[j]] = vals[j]
			}
			curr = goqu.And(eq, lt)
		}
		if expr == nil {
			expr = curr
		} else {
			expr = goqu.Or(expr, curr)
		}
	}
	return expr, nil
}

// * scanRows сканирует pgx.Rows в срез структур T.
// * Тип T должен быть struct, поля должны иметь тэг `db:"col"` или совпадать по имени с колонками.
func scanRows[T any](rows pgx.Rows) ([]T, error) {
	var zero T
	typ := reflect.TypeOf(zero)
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("T must be a struct")
	}

	//* Получаем описание колонок из результата
	fds := rows.FieldDescriptions()
	//* Колонка -> индекс поля в struct
	colIdx := make(map[string]int, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		name := sf.Tag.Get("db")
		if name == "" {
			name = sf.Name
		}
		colIdx[name] = i
	}

	var result []T
	for rows.Next() {
		//* Создаём новый экземпляр T
		ptr := reflect.New(typ)
		val := ptr.Elem()
		//* Подготавливаем слайс интерфейсов для Scan
		dests := make([]interface{}, len(fds))
		for i, fd := range fds {
			if idx, ok := colIdx[fd.Name]; ok {
				dests[i] = val.Field(idx).Addr().Interface()
			} else {
				var discard interface{}
				dests[i] = &discard
			}
		}
		if err := rows.Scan(dests...); err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}
		result = append(result, ptr.Elem().Interface().(T))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return result, nil
}

// * extractValuesReflect извлекает значения полей структуры item в том порядке,
// * в котором указаны cols. Используется для формирования курсора.
func extractValuesReflect[T any](item T, cols []string) ([]interface{}, error) {
	typ := reflect.TypeOf(item)
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("item must be a struct")
	}
	val := reflect.ValueOf(item)
	//* Карта имя колонки -> индекс поля
	fieldMap := make(map[string]int, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		name := sf.Tag.Get("db")
		if name == "" {
			name = sf.Name
		}
		fieldMap[name] = i
	}
	//* Собираем значения в порядке cols
	vals := make([]interface{}, len(cols))
	for i, col := range cols {
		idx, ok := fieldMap[col]
		if !ok {
			return nil, fmt.Errorf("column %s not found in struct", col)
		}
		vals[i] = val.Field(idx).Interface()
	}
	return vals, nil
}

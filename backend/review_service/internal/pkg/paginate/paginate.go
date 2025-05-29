// * Пакет paginate предоставляет универсальный механизм курсорной пагинации для баз данных PostgreSQL.
// * Он использует goqu для построения SQL-запросов и pgx для выполнения запросов, применяя зашифрованные курсоры
// * для навигации по наборам результатов без использования смещений.
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
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// * init регистрирует типы, используемые в курсорах, в пакете encoding/gob для сериализации.
func init() {
	gob.Register(time.Time{})          // * Регистрация time.Time для курсоров с отметкой времени
	gob.Register(uuid.UUID{})          // * Регистрация uuid.UUID для UUID-курсоров
	gob.Register([]byte{})             // * Регистрация []byte для бинарных данных
	gob.Register(pgtype.Timestamptz{}) // * Регистрация типа pgx timestamptz
	gob.Register(pgtype.UUID{})        // * Регистрация типа pgx UUID
}

// * TokenCoder определяет интерфейс для шифрования и расшифровки токенов курсора.
// * Реализации должны обеспечивать безопасное кодирование токенов для хранения на стороне клиента.
type TokenCoder interface {
	// * Encrypt принимает исходные данные в виде среза байт и возвращает строковый токен или ошибку.
	Encrypt(data []byte) (string, error)
	// * Decrypt дешифрует токен и возвращает исходные данные или ошибку при неудаче.
	Decrypt(token string) ([]byte, error)
}

// * Paginator[T] управляет состоянием и выполнением курсорной пагинации для типа T.
// * Используйте WithDataset, WithColumns и WithLimit для настройки запроса перед вызовом Paginate.
// * Paginate возвращает срез результатов типа T и зашифрованный токен для следующей страницы, если таковая имеется.
type Paginator[T any] struct {
	conn    *pgxpool.Pool       // * conn — пул соединений PostgreSQL
	dialect goqu.DialectWrapper // * dialect — SQL-диалект для построения запросов
	coder   TokenCoder          // * coder — шифратор/дешифратор токенов курсора
	limit   uint                // * limit — максимальное число записей на страницу
	cols    []string            // * cols — список колонок для сортировки и формирования курсора
	ds      *goqu.SelectDataset // * ds — базовый запрос (dataset) для пагинации
}

// * DefaultLimit задаёт размер страницы по умолчанию, если не указан через WithLimit.
const DefaultLimit = 20

// * NewPaginator создаёт новый Paginator.
// * Параметры:
// *   - conn: инициализированный pgxpool.Pool для подключения к БД (обязателен).
// *   - dialect: диалект goqu для построения SQL-запросов.
// *   - coder: реализация TokenCoder для операций с токенами (обязателен).
// * Возвращает:
// *   - *Paginator[T]: настроенный Paginator с лимитом по умолчанию.
// *   - error: если conn или coder равны nil.
func NewPaginator[T any](conn *pgxpool.Pool, dialect goqu.DialectWrapper, coder TokenCoder) (*Paginator[T], error) {
	if conn == nil {
		return nil, errors.New("db connection is nil")
	}
	if coder == nil {
		return nil, errors.New("coder is nil")
	}
	return &Paginator[T]{
		conn:    conn,
		dialect: dialect,
		coder:   coder,
		limit:   DefaultLimit,
	}, nil
}

// * WithDataset задаёт базовый SelectDataset для пагинации.
// * Использовать перед вызовом Paginate для установки контекста запроса.
// * Возвращает тот же Paginator для цепочки вызовов.
func (p *Paginator[T]) WithDataset(ds *goqu.SelectDataset) *Paginator[T] {
	p.ds = ds
	return p
}

// * WithLimit задаёт свой размер страницы.
// * Если limit равен нулю, текущий лимит сохраняется.
// * Возвращает тот же Paginator для цепочки.
func (p *Paginator[T]) WithLimit(limit uint) *Paginator[T] {
	if limit > 0 {
		p.limit = limit
	}
	return p
}

// * WithColumns настраивает колонки для сортировки и формирования курсоров.
// * Порядок колонок определяет лексикографический порядок при пагинации.
// * Возвращает тот же Paginator для цепочки.
func (p *Paginator[T]) WithColumns(cols ...string) *Paginator[T] {
	p.cols = append([]string(nil), cols...)
	return p
}

// * Paginate выполняет запрос для получения одной страницы результатов.
// * Параметры:
// *   - ctx: контекст выполнения запроса.
// *   - token: зашифрованный курсор предыдущей страницы; пустая строка для первой страницы.
// * Возвращает:
// *   - []T: срез результатов размером до limit.
// *   - string: токен для следующей страницы или пустую строку, если страниц больше нет.
// *   - error: ошибка при построении SQL, выполнении запроса или сканировании строк.
func (p *Paginator[T]) Paginate(ctx context.Context, token string) ([]T, string, error) {
	if p.limit == 0 {
		return nil, "", errors.New("limit must be > 0")
	}
	if p.ds == nil {
		return nil, "", errors.New("dataset not provided")
	}
	if len(p.cols) == 0 {
		return nil, "", errors.New("no columns specified for cursor")
	}

	// * Формируем выражения ORDER BY по убыванию для каждой колонки.
	orderExprs := make([]exp.OrderedExpression, len(p.cols))
	for i, c := range p.cols {
		orderExprs[i] = goqu.C(c).Desc()
	}
	ds := p.ds.Order(orderExprs...).Limit(p.limit + 1)

	// * Если передан токен, расшифровываем и применяем условие курсора.
	if token != "" {
		raw, err := p.coder.Decrypt(token)
		if err != nil {
			return nil, "", fmt.Errorf("token decrypt error: %w", err)
		}
		var vals []interface{}
		if err := gob.NewDecoder(bytes.NewReader(raw)).Decode(&vals); err != nil {
			return nil, "", fmt.Errorf("cursor decoding error: %w", err)
		}
		expr, err := cursorConditionExpr(p.cols, vals)
		if err != nil {
			return nil, "", err
		}
		ds = ds.Where(expr)
	}

	// * Генерируем SQL и выполняем запрос.
	sqlStr, args, err := ds.ToSQL()
	if err != nil {
		return nil, "", fmt.Errorf("SQL build error: %w", err)
	}
	rows, err := p.conn.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, "", fmt.Errorf("query execution error: %w", err)
	}
	defer rows.Close()

	// * Сканируем строки в срез структур T.
	items, err := scanRows[T](rows)
	if err != nil {
		return nil, "", err
	}

	// * Если получено больше элементов, чем limit, генерируем токен для следующей страницы.
	if len(items) > int(p.limit) {
		vals, err := extractValuesReflect(items[p.limit-1], p.cols)
		if err != nil {
			return nil, "", err
		}
		var buf bytes.Buffer
		if err := gob.NewEncoder(&buf).Encode(vals); err != nil {
			return nil, "", err
		}
		tok, err := p.coder.Encrypt(buf.Bytes())
		if err != nil {
			return nil, "", err
		}
		return items[:p.limit], tok, nil
	}
	return items, "", nil
}

// * cursorConditionExpr строит лексикографическое условие WHERE для фильтрации по курсору.
// * Возвращает выражение вида (col1<val1) OR (col1=val1 AND col2<val2) OR ...
func cursorConditionExpr(cols []string, vals []interface{}) (goqu.Expression, error) {
	if len(cols) != len(vals) {
		return nil, errors.New("cols and vals length mismatch")
	}
	var expr goqu.Expression
	for i := range cols {
		lt := goqu.Ex{cols[i]: goqu.Op{"lt": vals[i]}}
		var curr goqu.Expression
		if i == 0 {
			curr = lt
		} else {
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

// * scanRows считывает все строки из pgx.Rows и маппит их в срез структур типа T.
// * T должен быть структурой с полями, помеченными тегом `db:"column"` или совпадающими по имени с названием колонки.
// * Возвращает ошибку, если сканирование не удалось или T не является структурой.
func scanRows[T any](rows pgx.Rows) ([]T, error) {
	var zero T
	typ := reflect.TypeOf(zero)
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("T must be a struct")
	}

	// * Строим карту соответствия названий колонок и индексов полей структуры.
	fds := rows.FieldDescriptions()
	colIdx := make(map[string]int)
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
		ptr := reflect.New(typ)
		val := ptr.Elem()
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
			return nil, fmt.Errorf("row scan error: %w", err)
		}
		result = append(result, ptr.Elem().Interface().(T))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return result, nil
}

// * extractValuesReflect извлекает значения колонок курсора из экземпляра структуры через reflection.
// * Параметры:
// *   - item: экземпляр структуры типа T.
// *   - cols: список имён колонок для извлечения.
// * Возвращает:
// *   - []interface{}: значения в порядке, соответствующем cols.
// *   - error: если item не структура или колонка не найдена.
func extractValuesReflect[T any](item T, cols []string) ([]interface{}, error) {
	typ := reflect.TypeOf(item)
	if typ.Kind() != reflect.Struct {
		return nil, errors.New("item must be a struct")
	}
	val := reflect.ValueOf(item)

	// * Строим карту соответствия названий колонок и индексов полей структуры.
	fieldMap := make(map[string]int)
	for i := 0; i < typ.NumField(); i++ {
		sf := typ.Field(i)
		name := sf.Tag.Get("db")
		if name == "" {
			name = sf.Name
		}
		fieldMap[name] = i
	}

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

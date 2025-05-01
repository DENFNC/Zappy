// * Пакет paginate предоставляет утилиты для шифрования и дешифрования
// * токенов пагинации с использованием AES-GCM и кодирования в URL-безопасный Base64.
// *
// * Он определяет тип Encryptor, который инкапсулирует симметричный ключ
// * и источник случайных чисел, предлагая методы Encrypt и Decrypt
// * для работы с произвольными байтовыми срезами.
package paginate

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// * Encryptor выполняет шифрование и дешифрование данных
// * с помощью AES в режиме Galois/Counter Mode (GCM).
// * Он хранит симметричный ключ и источник случайных чисел для генерации nonce.
type Encryptor struct {
	//* key содержит ключ AES; допустимые длины — 16, 24 или 32 байта,
	//* соответствующие AES-128, AES-192 и AES-256 соответственно.
	key []byte
	//* rand — источник случайных чисел для генерации nonce;
	//* если nil, используется crypto/rand.Reader.
	rand io.Reader
}

// * NewEncryptor создаёт новый Encryptor с заданным ключом и источником случайных чисел.
// * Если randReader равен nil, по умолчанию используется crypto/rand.Reader.
// * Возвращает ошибку, если длина ключа некорректна.
// *
// * Параметры:
// *   - key: ключ AES (16, 24 или 32 байта)
// *   - randReader: источник случайных байт или nil для использования по умолчанию
// *
// * Возвращает:
// *   - *Encryptor: инициализированный шифровальщик
// *   - error: ненулевая при недопустимой длине ключа
func NewEncryptor(key []byte, randReader io.Reader) (*Encryptor, error) {
	keyLen := len(key)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return nil, fmt.Errorf("недопустимая длина ключа: %d", keyLen)
	}

	if randReader == nil {
		randReader = rand.Reader
	}
	return &Encryptor{key: key, rand: randReader}, nil
}

// *Encrypt шифрует переданные данные data с помощью AES-GCM.
// *Сначала создаётся случайный nonce, затем результат (nonce + шифр текст)
// *кодируется в URL-безопасный Base64.
// *
// *Параметры:
// *  - data: байты исходных данных для шифрования
// *
// *Возвращает:
// *  - string: закодированный Base64 результат со встроенным nonce
// *  - error: ненулевая при ошибках инициализации или шифрования
func (e *Encryptor) Encrypt(data []byte) (string, error) {
	//* Инициализация блока AES
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("ошибка инициализации шифра: %w", err)
	}
	//* Создание GCM-инстанса
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("ошибка инициализации GCM: %w", err)
	}
	//* Выделение места под nonce
	nonce := make([]byte, gcm.NonceSize())
	//* Заполнение nonce случайными данными
	if _, err := io.ReadFull(e.rand, nonce); err != nil {
		return "", fmt.Errorf("ошибка генерации nonce: %w", err)
	}
	//* Шифрование: результат содержит nonce + шифр текст
	cipherText := gcm.Seal(nonce, nonce, data, nil)
	//* Кодирование в URL-безопасный Base64
	return base64.URLEncoding.EncodeToString(cipherText), nil
}

// * Decrypt дешифрует токен, полученный из Encrypt.
// * Сначала декодируется Base64, затем извлекается nonce и шифр текст,
// * далее выполняется дешифрование и проверка целостности.
// *
// * Параметры:
// *   - token: строка URL-безопасного Base64 с префиксным nonce
// *
// * Возвращает:
// *   - []byte: восстановленные исходные данные
// *   - error: ненулевая при ошибках декодирования, некорректной длине токена
// *     или при сбое аутентификации/дешифрования
func (e *Encryptor) Decrypt(token string) ([]byte, error) {
	//* Декодирование Base64
	raw, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования токена: %w", err)
	}
	//* Инициализация блока AES
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации шифра: %w", err)
	}
	//* Создание GCM-инстанса
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации GCM: %w", err)
	}
	nonceSize := gcm.NonceSize()
	//* Проверка длины полученных данных
	if len(raw) < nonceSize {
		return nil, errors.New("некорректная длина токена")
	}
	//* Разделение на nonce и шифр текст
	nonce, ciphertext := raw[:nonceSize], raw[nonceSize:]
	//* Дешифрование и проверка целостности
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("ошибка дешифрования: %w", err)
	}
	return plaintext, nil
}

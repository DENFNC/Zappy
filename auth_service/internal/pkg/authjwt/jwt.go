// Package vaulttoken предоставляет функции для создания JWT-токенов, подписанных с помощью метода,
// основанного на публичном ключе из внешнего хранилища (vault). Это позволяет удобно генерировать
// безопасные токены с заданными регистрационными данными.
package vaulttoken

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Generate создает новый JWT-токен с указанными регистрационными данными и возвращает его в виде строки.
// Токен подписывается с помощью настроенного метода подписи, использующего публичный ключ, получаемый через vault.
//
// Параметры:
//   - vault: объект, реализующий интерфейс PublicKeyProvider, который обеспечивает доступ к публичному ключу для подписи.
//   - name: имя, используемое для определения конкретного метода подписи.
//   - iss: строка-издатель, которая указывается в регистрационных данных токена.
//   - expires: время действия токена; время истечения задается как текущая дата плюс expires, умноженное на time.Hour.
//
// Возвращаемые значения:
//   - строка с подписанным JWT-токеном.
//   - ошибка, если в процессе создания или подписания токена возникли проблемы.
func Generate(
	vault VaultKMS,
	iss, keyName string,
	expires time.Duration,
) (string, error) {
	// Создаем новый метод подписи, используя vault и имя
	newMethod := NewSigningMethodVaultPS256(vault, keyName)

	// Формируем стандартные зарегистерованные клеймы токена
	claims := jwt.NewWithClaims(
		newMethod,
		jwt.RegisteredClaims{
			Issuer:    iss,
			Subject:   "TODO",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expires * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	)

	// Подписываем токен, используя vault как ключ
	tokenStr, err := claims.SignedString(vault)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

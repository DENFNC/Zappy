// Package authjwt реализует алгоритмы подписи и верификации токенов JWT с использованием Hashicorp Vault.
// Реализована подпись с использованием алгоритма RSA-PSS (PS256), где приватный ключ хранится в Vault,
// а публичный ключ извлекается посредством интерфейса.
package authjwt

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
)

const (
	emptyValue = 0
)

// PublicKeyProvider определяет интерфейс для получения публичного RSA ключа.
// Реализации должны предоставлять метод FetchPublicKey, который возвращает
// публичный ключ и ошибку, если операция неудачна.
type PublicKeyProvider interface {
	// FetchPublicKey извлекает публичный RSA ключ.
	FetchPublicKey() (*rsa.PublicKey, error)
}

// VaultKey расширяет PublicKeyProvider и предоставляет методы для доступа
// к данным, необходимым для взаимодействия с Vault, таким как имя ключа и клиент.
type VaultKey interface {
	PublicKeyProvider
	GetKeyName() string
	GetClient() *api.Client
}

// SigningMethodVaultPS256 реализует метод подписи, основанный на алгоритме PS256 (RSA-PSS с SHA256)
// с использованием Vault для выполнения криптографических операций.
type SigningMethodVaultPS256 struct {
	Vault PublicKeyProvider
	name  string
}

func NewSigningMethodVaultPS256(
	vault PublicKeyProvider,
	name string,
) *SigningMethodVaultPS256 {
	return &SigningMethodVaultPS256{
		Vault: vault,
		name:  name,
	}
}

// Alg возвращает имя алгоритма подписи, используемого в SigningMethodVaultPS256.
func (m *SigningMethodVaultPS256) Alg() string {
	return m.name
}

// Sign подписывает строку signingString, используя приватный ключ, хранящийся в Vault.
// Параметр key должен реализовывать интерфейс VaultKey, который предоставляет необходимые данные для работы с Vault.
//
// Возвращает:
//   - []byte: подпись в виде байтов, полученная из Vault.
//   - error: подробную ошибку, если операция не выполнена.
func (m *SigningMethodVaultPS256) Sign(signingString string, key interface{}) ([]byte, error) {
	vaultKey, ok := key.(VaultKey)
	if !ok {
		return nil, fmt.Errorf("signing error: invalid key type; expected object implementing VaultKey")
	}

	digest := generateDigest(signingString)
	payload := createVaultSignData(digest)
	secret, err := performVaultSignOperation(vaultKey, payload)
	if err != nil {
		return nil, err
	}

	sigRaw, err := extractSignatureFromVaultResponse(secret)
	if err != nil {
		return nil, err
	}

	return sigRaw, nil
}

// generateDigest вычисляет SHA256-хэш для заданной строки.
//
// Параметры:
//   - signingString: строка, которую необходимо хэшировать.
//
// Возвращает:
//   - []byte: результат хэширования.
func generateDigest(signingString string) []byte {
	h := sha256.New()
	h.Write([]byte(signingString))
	return h.Sum(nil)
}

// createVaultSignData подготавливает данные для запроса в Vault на основе полученного хэша.
// Данные включают Base64-кодированный хэш и параметры алгоритмов подписи.
//
// Параметры:
//   - digest: SHA256-хэш исходной строки.
//
// Возвращает:
//   - map[string]any: мапа с параметрами, соответствующая требованиям Vault API.
func createVaultSignData(digest []byte) map[string]any {
	inputB64 := base64.StdEncoding.EncodeToString(digest)
	return map[string]any{
		"input":               inputB64,
		"signature_algorithm": "pss",
		"hash_algorithm":      "sha2-256",
		"prehashed":           true,
	}
}

// performVaultSignOperation выполняет запрос к Vault для выполнения операции подписи.
// Формируется путь вида "transit/sign/<KeyName>" и отправляется запрос с данными payload.
//
// Параметры:
//   - vaultKey: объект, реализующий интерфейс VaultKey, предоставляющий имя ключа и клиент Vault.
//   - data: данные для запроса в формате map[string]any.
//
// Возвращает:
//   - *api.Secret: результат операции от Vault.
//   - error: подробная ошибка, если запрос не выполнен или получен некорректный ответ.
func performVaultSignOperation(vaultKey VaultKey, data map[string]any) (*api.Secret, error) {
	signPath := fmt.Sprintf("transit/sign/%s", vaultKey.GetKeyName())
	secret, err := vaultKey.GetClient().Logical().Write(signPath, data)
	if err != nil {
		return nil, fmt.Errorf("vault sign error: failed to write to %q: %w", signPath, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("vault sign error: received empty response from %q", signPath)
	}
	if len(secret.Warnings) != emptyValue {
		return nil, fmt.Errorf("vault sign error: warning received: %s", secret.Warnings[0])
	}
	return secret, nil
}

// extractSignatureFromVaultResponse извлекает строку подписи из ответа Vault,
// разбивает её по разделителю ":" и декодирует итоговую часть из Base64.
//
// Параметры:
//   - secret: объект *api.Secret, полученный из Vault.
//
// Возвращает:
//   - []byte: декодированная подпись.
//   - error: подробное описание ошибки, если формат подписи неверный или декодирование не удалось.
func extractSignatureFromVaultResponse(secret *api.Secret) ([]byte, error) {
	vaultSignature, ok := secret.Data["signature"].(string)
	if !ok {
		return nil, fmt.Errorf("vault sign error: missing 'signature' field in response")
	}
	parts := strings.SplitN(vaultSignature, ":", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("vault sign error: unexpected signature format %q", vaultSignature)
	}
	sigRaw, err := base64.StdEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, fmt.Errorf("vault sign error: failed to decode signature from Base64: %w", err)
	}
	return sigRaw, nil
}

// Verify проверяет корректность подписи для заданной строки signingString,
// используя публичный ключ, полученный посредством интерфейса PublicKeyProvider.
// Параметр key должен реализовывать PublicKeyProvider.
//
// Параметры:
//   - signingString: исходная строка, для которой проверяется подпись.
//   - sig: подпись в виде байтов, которую необходимо проверить.
//   - key: объект, реализующий PublicKeyProvider, для получения публичного ключа.
//
// Возвращает:
//   - error: nil, если подпись корректна; в противном случае – подробное описание ошибки.
func (m *SigningMethodVaultPS256) Verify(signingString string, sig []byte, key interface{}) error {
	vsk, ok := key.(PublicKeyProvider)
	if !ok {
		return fmt.Errorf("verify error: invalid key type; expected object implementing PublicKeyProvider")
	}

	pubKey, err := vsk.FetchPublicKey()
	if err != nil {
		return fmt.Errorf("verify error: failed to fetch public key from vault: %w", err)
	}

	h := sha256.New()
	h.Write([]byte(signingString))
	digest := h.Sum(nil)

	if err := rsa.VerifyPSS(pubKey, crypto.SHA256, digest, sig, &rsa.PSSOptions{
		SaltLength: rsa.PSSSaltLengthAuto,
		Hash:       crypto.SHA256,
	}); err != nil {
		return fmt.Errorf("verify error: RSA-PSS verification failed: %w", err)
	}
	return nil
}

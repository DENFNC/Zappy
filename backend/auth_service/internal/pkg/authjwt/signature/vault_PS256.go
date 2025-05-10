// Package vaulttoken реализует алгоритмы подписи и верификации токенов JWT с использованием Hashicorp Vault.
// Реализована подпись с использованием алгоритма RSA-PSS (PS256), где приватный ключ хранится в Vault,
// а публичный ключ извлекается посредством интерфейса.
package signature

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/hashicorp/vault/api"
)

const (
	methodVaultPS256 = "PS256"
)

type VaultKMS interface {
	Read(keyName string) (*api.Secret, error)
	Sign(keyName string, data map[string]any) (*api.Secret, error)
}

// SigningMethodVaultPS256 реализует метод подписи, основанный на алгоритме PS256(RSA-PSS с SHA256)
// с использованием Vault для выполнения криптографических операций.
type SigningMethodVaultPS256 struct {
	Vault   VaultKMS
	keyName string
	name    string
}

func NewSigningMethodVaultPS256(
	vault VaultKMS,
	keyName string,
) *SigningMethodVaultPS256 {
	return &SigningMethodVaultPS256{
		Vault:   vault,
		keyName: keyName,
		name:    methodVaultPS256,
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
	digest := generateDigest(signingString)
	payload := createVaultSignData(digest)
	secret, err := m.Vault.Sign(m.keyName, payload)
	if err != nil {
		return nil, err
	}

	sigRaw, err := extractSignatureFromVaultResponse(secret)
	if err != nil {
		return nil, err
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
	pubKey, err := m.fetchPublicKey()
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

// fetchPublicKey получает публичный RSA ключ из Vault для данного ключа подписания.
// Данный метод декомпозирован на несколько вспомогательных функций, каждая из которых
// отвечает за отдельный этап обработки:
//  1. Чтение информации о ключе из Vault (readKeySecret).
//  2. Извлечение PEM-представления публичного ключа (extractPublicKeyPEM).
//  3. Декодирование и парсинг PEM в *rsa.PublicKey (parseRSAPublicKeyFromPEM).
//
// Возвращает:
//   - *rsa.PublicKey: указатель на публичный RSA ключ, если процесс извлечения и парсинга прошел успешно.
//   - error: ошибка, описывающая возникшую проблему на любом из этапов получения или обработки данных.
func (m *SigningMethodVaultPS256) fetchPublicKey() (*rsa.PublicKey, error) {
	secret, err := m.Vault.Read(m.keyName)
	if err != nil {
		return nil, err
	}

	pubPem, err := extractPublicKeyPEM(secret, m.keyName)
	if err != nil {
		return nil, err
	}

	return parseRSAPublicKeyFromPEM(pubPem)
}

// extractPublicKeyPEM извлекает строку с PEM-представлением публичного ключа из ответа Vault.
// Для этого из данных ключей выбирается информация о последней версии.
func extractPublicKeyPEM(secret *api.Secret, keyName string) (string, error) {
	keysAny, ok := secret.Data["keys"]
	if !ok {
		return "", fmt.Errorf("no 'keys' field in response from vault for key: %s", keyName)
	}
	keysMap, ok := keysAny.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("cannot parse 'keys' field as map for key: %s", keyName)
	}

	versionDataAny, ok := keysMap["latest_version"]
	if !ok {
		return "", fmt.Errorf("failed to find latest key version for key: %s", keyName)
	}

	versionData, ok := versionDataAny.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid data format for key version for key: %s", keyName)
	}

	pubPem, ok := versionData["public_key"].(string)
	if !ok {
		return "", fmt.Errorf("no 'public_key' field for latest key version of key: %s", keyName)
	}

	return pubPem, nil
}

// parseRSAPublicKeyFromPEM декодирует PEM-представление публичного ключа
// и парсит его в объект *rsa.PublicKey.
func parseRSAPublicKeyFromPEM(pubPem string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPem))
	if block == nil {
		return nil, fmt.Errorf("failed to PEM-decode public key")
	}

	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPub, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not an RSA key")
	}
	return rsaPub, nil
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

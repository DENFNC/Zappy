// Package main демонстрирует инициализацию клиента Vault для подписания ключей,
// используя API Hashicorp Vault.
package vault

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/api"
)

// VaultSignKey инкапсулирует клиента Vault и имя ключа, используемого для подписания.
type VaultSignKey struct {
	Client  *api.Client
	KeyName string
}

// NewVaultSignKey создает и возвращает новый экземпляр VaultSignKey.
// Функция проверяет, что переданные параметры не пустые, и инициализирует клиента Vault.
//
// Параметры:
//   - name: не пустой идентификатор для ключа подписания.
//   - addr: адрес сервера Vault (например, "http://127.0.0.1:8200").
//   - token: валидный токен для доступа к Vault.
//
// Возвращает:
//   - *VaultSignKey: указатель на инициализированный объект VaultSignKey.
//   - error: ошибку, описывающую причину сбоя, если что-то пошло не так.
func NewVaultSignKey(name, addr, token string) (*VaultSignKey, error) {
	if name == "" {
		return nil, errors.New("key name must not be empty")
	}

	newClient, err := initVault(addr, token)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize Vault client: %w", err)
	}

	return &VaultSignKey{
		Client:  newClient,
		KeyName: name,
	}, nil
}

// GetKeyName возвращает имя ключа для подписания.
func (vsk *VaultSignKey) GetKeyName() string {
	return vsk.KeyName
}

// GetClient возвращает клиента Vault.
func (vsk *VaultSignKey) GetClient() *api.Client {
	return vsk.Client
}

// initVault инициализирует клиента Vault, принимая адрес и токен для подключения.
// Проверяются входные параметры, чтобы они не были пустыми. Затем создается клиент
// с использованием стандартной конфигурации Vault API.
//
// Параметры:
//   - addr: адрес сервера Vault; не должен быть пустым.
//   - token: токен для доступа к Vault; не должен быть пустым.
//
// Возвращает:
//   - *api.Client: указатель на настроенный клиент Vault.
//   - error: ошибку с описанием, если валидация параметров не прошла или произошла ошибка при создании клиента.
func initVault(addr, token string) (*api.Client, error) {
	if addr == "" {
		return nil, errors.New("vault address cannot be empty")
	}
	if token == "" {
		return nil, errors.New("vault token cannot be empty")
	}

	cfg := api.DefaultConfig()
	cfg.Address = addr

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating client: %w", err)
	}

	client.SetToken(token)
	return client, nil
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
func (vsk *VaultSignKey) FetchPublicKey() (*rsa.PublicKey, error) {
	secret, err := vsk.readKeySecret()
	if err != nil {
		return nil, err
	}

	pubPem, err := extractPublicKeyPEM(secret, vsk.KeyName)
	if err != nil {
		return nil, err
	}

	return parseRSAPublicKeyFromPEM(pubPem)
}

// readKeySecret выполняет запрос к Vault по пути "transit/keys/<KeyName>"
// для получения информации о ключе.
func (vsk *VaultSignKey) readKeySecret() (*api.Secret, error) {
	path := fmt.Sprintf("transit/keys/%s", vsk.KeyName)
	secret, err := vsk.Client.Logical().Read(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key info from vault: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("key info not found in vault: %s", path)
	}
	return secret, nil
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

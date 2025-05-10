package signature

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

const (
	methodVaultEdDSA = "EdDSA"
)

// type VaultKMS interface {
// 	Read(keyName string) (*api.Secret, error)
// 	Sign(keyName string, data map[string]any) (*api.Secret, error)
// }

type SigningMethodVaultEdDSA struct {
	Vault   VaultKMS
	keyName string
	name    string
}

func NewSigningMethodVaultEdDSA(
	vault VaultKMS,
	keyName string,
) *SigningMethodVaultEdDSA {
	return &SigningMethodVaultEdDSA{
		Vault:   vault,
		keyName: keyName,
		name:    methodVaultEdDSA,
	}
}

func (m *SigningMethodVaultEdDSA) Alg() string {
	return m.name
}

func (m *SigningMethodVaultEdDSA) Sign(signingString string, key interface{}) ([]byte, error) {
	payload := map[string]any{
		"input": base64.StdEncoding.EncodeToString([]byte(signingString)),
	}
	secret, err := m.Vault.Sign(m.keyName, payload)
	if err != nil {
		return nil, err
	}

	return extractSignatureFromVaultResponse(secret)
}

func (m *SigningMethodVaultEdDSA) Verify(signingString string, sig []byte, key interface{}) error {
	pubKey, err := m.fetchPublicKey()
	if err != nil {
		return fmt.Errorf("verify error: failed to fetch public key from vault: %w", err)
	}

	if !ed25519.Verify(pubKey, []byte(signingString), sig) {
		return fmt.Errorf("verify error: EdDSA verification failed")
	}
	return nil
}

func (m *SigningMethodVaultEdDSA) fetchPublicKey() (ed25519.PublicKey, error) {
	secret, err := m.Vault.Read(m.keyName)
	if err != nil {
		return nil, err
	}

	pubPem, err := extractPublicKeyPEM(secret, m.keyName)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode([]byte(pubPem))
	if block == nil {
		return nil, fmt.Errorf("failed to PEM-decode public key")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	edPub, ok := pubKey.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not Ed25519")
	}
	return edPub, nil
}

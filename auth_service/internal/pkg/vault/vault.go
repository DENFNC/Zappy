package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
)

const (
	emptyValue = 0
)

type Vault struct {
	client *api.Client
}

func New(addr, authToken string) (*Vault, error) {
	cfg := api.DefaultConfig()
	cfg.Address = addr

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}

	client.SetToken(authToken)

	return &Vault{
		client: client,
	}, nil
}

func (v *Vault) Read(keyName string) (*api.Secret, error) {
	pathRead := fmt.Sprintf("transit/keys/%s", keyName)
	secret, err := v.client.Logical().Read(pathRead)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (v *Vault) Sign(keyName string, data map[string]any) (*api.Secret, error) {
	pathSign := fmt.Sprintf("transit/sign/%s", keyName)
	secret, err := v.client.Logical().Write(pathSign, data)
	if err != nil {
		return nil, fmt.Errorf("vault sign error: failed to write to %q: %w", pathSign, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("vault sign error: received empty response from %q", pathSign)
	}
	if len(secret.Warnings) != emptyValue {
		return nil, fmt.Errorf("vault sign error: warning received: %s", secret.Warnings[0])
	}
	return secret, nil
}

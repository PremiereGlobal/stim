package vault

import (
	"errors"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/token"
)

type Vault struct {
	client *api.Client
}

func New() *Vault {
	v := &Vault{}
	return v
}

func (v *Vault) InitClient() error {
	// Initialize client
	client, err := api.NewClient(nil)
	if err != nil {
		return err
	}
	tokenHelper := token.InternalTokenHelper{}
	token, err := tokenHelper.Get()
	if err != nil {
		return err
	}
	client.SetToken(token)
	v.client = client

	return nil
}

func (v *Vault) Login() {

}

func (v *Vault) GetSecretKey(path string, key string) (string, error) {

	secret, err := v.client.Logical().Read(path)
	if err != nil {
		return "", err
	}

	// If we got back an empty response, fail
	if secret == nil {
		return "", errors.New("Could not find secret `" + path + "`")
	}

	// If the provided key doesn't exist, fail
	if secret.Data[key] == nil {
		return "", errors.New("Vault: Could not find key `" + key + "` for secret `" + path + "`")
	}

	return secret.Data[key].(string), nil
}

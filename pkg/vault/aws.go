package vault

import (
	"github.com/hashicorp/vault/api"

	"errors"
)

func (v *Vault) AWScredentials(account string, role string) (*api.Secret, error) {
	if account == "" {
		return nil, errors.New("Account not set")
	}
	if role == "" {
		return nil, errors.New("Role not set")
	}

	path := "/" + account + "/creds/" + role
	v.log.Debug("Getting AWS credentials via path: ", path)

	secret, err := v.GetSecret(path)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

package vault

import (
	"github.com/hashicorp/vault/api"
	"github.com/readytalk/stim/pkg/log"

	"errors"
)

func (v *Vault) AWScredentials(account string, role string, sts bool) (*api.Secret, error) {
	if account == "" {
		return nil, errors.New("Account not set")
	}
	if role == "" {
		return nil, errors.New("Role not set")
	}

	credType := "creds"
	if sts {
		credType = "sts"
	}
	path := "/" + account + "/" + credType + "/" + role
	log.Debug("Getting AWS credentials via path: ", path)

	secret, err := v.GetSecret(path)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

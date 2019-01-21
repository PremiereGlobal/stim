package vault

import (
	"errors"
)

// Pulls a single key from a secret path and returns the value as a string
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

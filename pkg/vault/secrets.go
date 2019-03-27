package vault

import (
	"github.com/hashicorp/vault/api"
	"path/filepath"
)

// Using Vaults Logical client:
// https://github.com/hashicorp/vault/blob/master/api/logical.go

// GetSecretKey takes a secret path and key and returns, if successful,
// the secret string present in that key.
func (v *Vault) GetSecretKey(path string, key string) (string, error) {

	secret, err := v.client.Logical().Read(path)
	if err != nil {
		return "", v.parseError(err)
	}

	// If we got back an empty response, fail
	if secret == nil {
		return "", v.newError("Could not find secret `" + path + "`")
	}

	// If the provided key doesn't exist, fail
	if secret.Data[key] == nil {
		return "", v.newError("Vault: Could not find key `" + key + "` for secret `" + path + "`")
	}

	return secret.Data[key].(string), nil
}

// GetSecretKeys takes a secret path and returns, if successful,
// a map of all the keys at that path.
func (v *Vault) GetSecretKeys(path string) (map[string]string, error) {

	secret, err := v.client.Logical().Read(path)
	if err != nil {
		return nil, v.parseError(err)
	}

	// If we got back an empty response, fail
	if secret == nil {
		return nil, v.newError("Could not find secret `" + path + "`")
	}

	// Loop through and get all the keys
	var secretList map[string]string
	secretList = make(map[string]string)
	for key, value := range secret.Data {
		secretList[key] = value.(string)
	}

	return secretList, nil
}

// ListSecrets takes a secret path and returns, if successful,
// a list of all child paths under that path.
func (v *Vault) ListSecrets(path string) ([]string, error) {

	secret, err := v.client.Logical().List(path)
	if err != nil {
		return nil, v.parseError(err)
	}

	// If we got back an empty response, fail
	if secret == nil {
		return nil, v.newError("Could not find secret `" + path + "`")
	}

	// Loop through and get all the keys
	var secretList []string
	for _, value := range secret.Data["keys"].([]interface{}) {
		secretList = append(secretList, filepath.Clean(value.(string)))
	}

	return secretList, nil
}

// GetSecret takes a secret path and returns the secret(s) in a Vault object
func (v *Vault) GetSecret(path string) (*api.Secret, error) {
	secret, err := v.client.Logical().Read(path)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

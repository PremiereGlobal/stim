package vault

import (
	"errors"
	// "github.com/davecgh/go-spew/spew"
	"path/filepath"
)

// GetSecretKey takes a secret path and key and returns, if successful,
// the secret string present in that key.
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

// GetSecretKeys takes a secret path and returns, if successful,
// a map of all the keys at that path.
func (v *Vault) GetSecretKeys(path string) (map[string]string, error) {

	secret, err := v.client.Logical().Read(path)
	if err != nil {
		return nil, err
	}

	// If we got back an empty response, fail
	if secret == nil {
		return nil, errors.New("Could not find secret `" + path + "`")
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
		return nil, err
	}

	// If we got back an empty response, fail
	if secret == nil {
		return nil, errors.New("Could not find secret `" + path + "`")
	}

	// Loop through and get all the keys
	var secretList []string
	for _, value := range secret.Data["keys"].([]interface{}) {
		secretList = append(secretList, filepath.Clean(value.(string)))
	}

	return secretList, nil
}

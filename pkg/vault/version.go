package vault

// Version returns the version of the Vault server
func (v *Vault) Version() (string, error) {

	result, err := v.GetHealth()
	if err != nil {
		return "", err
	}

	return result.Version, nil
}

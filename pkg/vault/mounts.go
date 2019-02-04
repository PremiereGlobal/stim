package vault

import (
// "github.com/davecgh/go-spew/spew"
)

// GetMounts retrieves a list of mounts
func (v *Vault) GetMounts(mountType string) ([]string, error) {

	mounts, err := v.client.Sys().ListMounts()
	if err != nil {
		return nil, err
	}

	var result []string
	for path, mountOutput := range mounts {
		if mountOutput.Type == mountType {
			result = append(result, path)
		}
	}

	return result, nil
}

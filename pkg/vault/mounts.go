package vault

import (
	"sort"
	"strings"
)

// GetMounts retrieves a list of mounts
// Will return a string array filtered with given type. Example 'aws'
func (v *Vault) GetMounts(mountType string) ([]string, error) {

	mounts, err := v.client.Sys().ListMounts()
	if err != nil {
		return nil, err
	}

	var result []string
	for path, mountOutput := range mounts {
		if mountOutput.Type == mountType {
			result = append(result, strings.TrimRight(path, "/"))
		}
	}

	sort.Strings(result)
	return result, nil
}

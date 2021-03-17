package vault

import (
	"github.com/PremiereGlobal/stim/pkg/utils"
)

func (v *Vault) Filter(paths []string, withCapabilities []string) ([]string, error) {
	if len(paths) == 0 {
		return []string{}, nil
	}

	opts := &CapabilitiesSelfOptions{
		Paths: paths,
	}

	results, err := v.CapabilitiesSelf(opts)
	if err != nil {
		return nil, err
	}

	if len(withCapabilities) == 0 {
		withCapabilities = []string{"list", "read"}
	}

	var filteredPaths []string
	for path, capabilities := range results.Data {
		for _, capability := range withCapabilities {
			if path == "capabilities" {
				continue
			}
			if utils.Contains(capabilities, capability) {
				filteredPaths = append(filteredPaths, path)
				break
			}
		}
	}

	return filteredPaths, nil
}

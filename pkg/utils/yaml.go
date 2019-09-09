package utils

import (
	"gopkg.in/yaml.v2"
)

// IsYaml verifies that a byte string is value YAML syntax
// On failure, returns an error with the reason
func IsYaml(s []byte) (bool, error) {
	var y map[string]interface{}
	err := yaml.Unmarshal(s, &y)
	return err == nil, err
}

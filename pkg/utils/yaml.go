package utils

import (
	"gopkg.in/yaml.v2"
)

func IsYaml(s []byte) bool {
	var y map[string]interface{}
	return yaml.Unmarshal(s, &y) == nil
}

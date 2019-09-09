package deploy

import (
	"encoding/json"
	"fmt"
)

// makeSecretConfig generates a secret config json string based on the instance configuration
func (d *Deploy) makeSecretConfig(instance *Instance) (string, error) {

	secretConfigString := ""

	if len(instance.EnvSpec.Secrets) > 0 {

		b, err := json.Marshal(instance.EnvSpec.Secrets)
		if err != nil {
			fmt.Println("error:", err)
		}

		secretConfigString = string(b)
	}

	return secretConfigString, nil
}

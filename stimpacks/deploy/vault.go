package deploy

import (
	"encoding/json"
)

// makeSecretConfig generates a secret config json string based on the instance configuration
func (d *Deploy) makeSecretConfig(instance *Instance) (string, error) {

	secretConfigString := ""

	if len(instance.Spec.Secrets) > 0 {

		b, err := json.Marshal(instance.Spec.Secrets)
		if err != nil {
			d.log.Fatal("Unable to create secret config: {}", err)
		}

		secretConfigString = string(b)
	}

	return secretConfigString, nil
}

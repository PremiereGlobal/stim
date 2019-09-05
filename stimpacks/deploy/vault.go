package deploy

import (
	"encoding/json"
	"fmt"
)

func (d *Deploy) makeSecretConfig(cluster *Cluster) (string, error) {

	secretConfigString := ""
	if len(cluster.EnvSpec.Secrets) > 0 {

		b, err := json.Marshal(cluster.EnvSpec.Secrets)
		if err != nil {
			fmt.Println("error:", err)
		}

		secretConfigString = string(b)
	}

	return secretConfigString, nil
}

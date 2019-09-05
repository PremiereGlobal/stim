package deploy

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"

	"github.com/PremiereGlobal/stim/pkg/utils"
	v2e "github.com/PremiereGlobal/vault-to-envs/pkg/vaulttoenvs"
)

type Config struct {
	filePath     string
	Globals      map[string]string      `yaml:"globals"`
	Environments map[string]Environment `yaml:"environments"`
}

type Environment struct {
	Clusters  map[string]Cluster
	Secrets   []v2e.SecretItem `yaml:"secrets"`
	Namespace string
}

type Cluster struct {
	Namespace      string           `yaml:"namespace"`
	ServiceAccount string           `yaml:"serviceAccount"`
	Secrets        []v2e.SecretItem `yaml:"secrets"`
}

//
// type SecretItem struct {
// 	SecretPath         string            `yaml:"path"`
// 	TTL                int               `json:"ttl"`
// 	Version            float64           `json:"version"`
// 	SecretMaps         map[string]string `json:"set"`
// 	SecretDataPath     string            // kv v2
// 	SecretMetadataPath string            // kv v2
// 	EffectiveVersion   int               // kv v2
// 	secretMapValues    map[string]string
// 	secret             *VaultApi.Secret
// 	mount              *VaultApi.MountOutput
// }

func (d *Deploy) ParseConfig() {

	d.config = Config{}

	configFile := d.stim.GetConfig("deploy.file")

	if configFile == "" {
		d.stim.Fatal(errors.New("Must provide deployment config file '--deploy-file'"))
	}

	_, err := os.Stat(configFile)
	if err != nil && !os.IsExist(err) {
		d.stim.Fatal(errors.New(fmt.Sprintf("No deployment config file exists at: %s", configFile)))
	}

	contentstring, err := ioutil.ReadFile(configFile)
	if err != nil {
		d.stim.Fatal(errors.New(fmt.Sprintf("Deployment config file could not be read: %v", err)))
	}

	if !utils.IsYaml(contentstring) {
		d.stim.Fatal(errors.New(fmt.Sprintf("Deployment config file is not valid YAML: %v", err)))
	}

	err = yaml.Unmarshal([]byte(contentstring), &d.config)
	if err != nil {
		d.stim.Fatal(errors.New(fmt.Sprintf("Error parsing deployment config %v", err)))
	}

	d.config.filePath = configFile

	d.ValidateConfig()

}

func (d *Deploy) ValidateConfig() {

	// Ensure the environment names don't have "All"
	for e := range d.config.Environments {
		for c := range d.config.Environments[e].Clusters {
			if strings.ToLower(c) == "all" {
				d.stim.Fatal(errors.New("Deployment config cannot have an cluster named 'All'. It is a reserved name for deploying to all clusters"))
			}
		}
	}
}

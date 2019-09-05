package deploy

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	// "github.com/imdario/mergo"
	// "github.com/davecgh/go-spew/spew"

	"github.com/PremiereGlobal/stim/pkg/utils"
	v2e "github.com/PremiereGlobal/vault-to-envs/pkg/vaulttoenvs"
)

const (
	DEFAULT_CONTAINER_REPO   = "premiereglobal/kube-vault-deploy"
	DEFAULT_CONTAINER_TAG    = "0.3.1"
	DEFAULT_DEPLOY_DIRECTORY = "./"
	DEFAULT_DEPLOY_SCRIPT    = "deploy.sh"
)

type Config struct {
	configFilePath string
	Deployment     Deployment              `yaml:"deployment"`
	Container      Container               `yaml:"container"`
	Global         Global                  `yaml:"global"`
	Environments   map[string]*Environment `yaml:"environments"`
}

type Deployment struct {
	Directory         string `yaml:"dir"`
	Script            string `yaml:"script"`
	fullDirectoryPath string
}

type Container struct {
	Repo string `yaml:"repo"`
	Tag  string `yaml:"tag"`
}

type Global struct {
	EnvSpec *EnvSpec `yaml:"envSpec"`
}

type EnvSpec struct {
	Namespace       string           `yaml:"namespace"`
	TillerNamespace string           `yaml:"tillerNamespace"`
	Secrets         []v2e.SecretItem `yaml:"secrets"`
	EnvironmentVars []EnvironmentVar `yaml:"env"`
}

type Environment struct {
	EnvSpec  *EnvSpec `yaml:"envSpec"`
	Clusters map[string]*Cluster
}

type Cluster struct {
	EnvSpec        *EnvSpec `yaml:"envSpec"`
	ServiceAccount string   `yaml:"serviceAccount"`
}

type EnvironmentVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
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
		d.stim.Fatal(errors.New(fmt.Sprintf("Deployment config file (%s) is not valid YAML: %v", configFile, err)))
	}

	err = yaml.Unmarshal([]byte(contentstring), &d.config)
	if err != nil {
		d.stim.Fatal(errors.New(fmt.Sprintf("Error parsing deployment config %v", err)))
	}

	d.config.configFilePath = configFile

	d.ProcessConfig()

}

func (d *Deploy) ProcessConfig() {

	// Set defaults
	setConfigDefault(&d.config.Container.Repo, DEFAULT_CONTAINER_REPO)
	setConfigDefault(&d.config.Container.Tag, DEFAULT_CONTAINER_TAG)
	setConfigDefault(&d.config.Deployment.Directory, DEFAULT_DEPLOY_DIRECTORY)
	setConfigDefault(&d.config.Deployment.Script, DEFAULT_DEPLOY_SCRIPT)

	for e := range d.config.Environments {
		for c := range d.config.Environments[e].Clusters {

			// Ensure the environment names don't have "All".  This is a reserved name for designating a deployment to all clusters in an environment
			if strings.ToLower(c) == "all" {
				d.stim.Fatal(errors.New("Deployment config cannot have an cluster named 'All'. It is a reserved name for deploying to all clusters"))
			}

			// Create our envSpecs if they doesn't exist so we can fill it with cascaded values
			if d.config.Global.EnvSpec == nil {
				d.config.Global.EnvSpec = &EnvSpec{}
			}
			if d.config.Environments[e].EnvSpec == nil {
				d.config.Environments[e].EnvSpec = &EnvSpec{}
			}
			if d.config.Environments[e].Clusters[c].EnvSpec == nil {
				d.config.Environments[e].Clusters[c].EnvSpec = &EnvSpec{}
			}

			// Makes these easier to access.  They are pointers so any updates will stick
			globalEnvSpec := d.config.Global.EnvSpec
			environmentEnvSpec := d.config.Environments[e].EnvSpec
			clusterEnvSpec := d.config.Environments[e].Clusters[c].EnvSpec

			// Perform our envSpec merges
			// This will roll down the envSpecs from global -> environment -> cluster where each lower config takes precedence
			// First merge the environment envSpec into the cluster envSpec
			// Then merge the global envSpec into the cluster envSpec
			setConfigDefault(&clusterEnvSpec.Namespace, environmentEnvSpec.Namespace)
			setConfigDefault(&clusterEnvSpec.Namespace, globalEnvSpec.Namespace)
			setConfigDefault(&clusterEnvSpec.TillerNamespace, environmentEnvSpec.TillerNamespace)
			setConfigDefault(&clusterEnvSpec.TillerNamespace, globalEnvSpec.TillerNamespace)

			// Now do the arrays by adding them to a list starting with the global values and moving down
			// This way the lowest values will take effect last (and thus taking precedence)
			var newClusterSecrets []v2e.SecretItem
			var newClusterEnvVars []EnvironmentVar
			if len(globalEnvSpec.Secrets) > 0 {
				newClusterSecrets = append(newClusterSecrets, globalEnvSpec.Secrets...)
			}
			if len(globalEnvSpec.EnvironmentVars) > 0 {
				newClusterEnvVars = append(newClusterEnvVars, globalEnvSpec.EnvironmentVars...)
			}
			if len(environmentEnvSpec.Secrets) > 0 {
				newClusterSecrets = append(newClusterSecrets, environmentEnvSpec.Secrets...)
			}
			if len(environmentEnvSpec.EnvironmentVars) > 0 {
				newClusterEnvVars = append(newClusterEnvVars, environmentEnvSpec.EnvironmentVars...)
			}
			if len(clusterEnvSpec.Secrets) > 0 {
				newClusterSecrets = append(newClusterSecrets, clusterEnvSpec.Secrets...)
			}
			if len(clusterEnvSpec.EnvironmentVars) > 0 {
				newClusterEnvVars = append(newClusterEnvVars, clusterEnvSpec.EnvironmentVars...)
			}

			clusterEnvSpec.Secrets = newClusterSecrets
			clusterEnvSpec.EnvironmentVars = newClusterEnvVars
		}
	}

	// Determine the full directory path
	d.config.Deployment.fullDirectoryPath = filepath.Join(filepath.Dir(d.config.configFilePath), d.config.Deployment.Directory)

	// spew.Dump(d.config)
	// spew.Dump("")

}

func setConfigDefault(value *string, def string) {
	if len(*value) == 0 {
		*value = def
	}
}

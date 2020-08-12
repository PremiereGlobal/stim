package deploy

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/PremiereGlobal/stim/pkg/utils"
	"github.com/PremiereGlobal/stim/stim"
	v2e "github.com/PremiereGlobal/vault-to-envs/pkg/vaulttoenvs"
	"gopkg.in/yaml.v2"
)

const (
	defaultContainerRepo   = "premiereglobal/kube-vault-deploy"
	defaultContainerTag    = "0.3.3"
	defaultDeployDirectory = "./"
	defaultDeployScript    = "deploy.sh"
	defaultConfigFile      = "./stim.deploy.yaml"
)

// Config is the root structure for the deployment configuration
type Config struct {
	configFilePath string
	Deployment     Deployment     `yaml:"deployment"`
	Global         Global         `yaml:"global"`
	Environments   []*Environment `yaml:"environments"`
	environmentMap map[string]int
}

// Deployment describes details about the deployment assets (directories, files, etc)
type Deployment struct {
	Directory         string    `yaml:"directory"`
	Script            string    `yaml:"script"`
	Container         Container `yaml:"container"`
	fullDirectoryPath string
}

// Container describes the container used for Docker deployments
type Container struct {
	Repo string `yaml:"repo"`
	Tag  string `yaml:"tag"`
}

// Global describes global environment specs
type Global struct {
	Spec *Spec `yaml:"spec"`
}

// Spec contains the spec of a given environment/instance
type Spec struct {
	Kubernetes            Kubernetes              `yaml:"kubernetes"`
	Secrets               []*v2e.SecretItem       `yaml:"secrets"`
	EnvironmentVars       []*EnvironmentVar       `yaml:"env"`
	AddConfirmationPrompt bool                    `yaml:"addConfirmationPrompt"`
	Tools                 map[string]stim.EnvTool `yaml:"tools"`
}

// Kubernetes describes the Kubernetes configuration to use
type Kubernetes struct {
	ServiceAccount string `yaml:"serviceAccount"`
	Cluster        string `yaml:"cluster"`
}

// Environment describes a deployment environment (i.e. dev, stage, prod, etc.)
type Environment struct {
	Name            string      `yaml:"name"`
	Spec            *Spec       `yaml:"spec"`
	Instances       []*Instance `yaml:"instances"`
	RemoveAllPrompt bool        `yaml:"removeAllPrompt"`
	instanceMap     map[string]int
}

// Instance describes an instance of a deployment within an environment (i.e. us-west-2 for env prod)
type Instance struct {
	Name string `yaml:"name"`
	Spec *Spec  `yaml:"spec"`
}

// EnvironmentVar describes a shell env var to be injected into the deployment environment
type EnvironmentVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// parseConfig opens the deployment config file and ensures it is valid
func (d *Deploy) parseConfig() {

	d.config = Config{}

	configFile := d.stim.ConfigGetString("deploy.file")

	if configFile == "" {
		setConfigDefault(&configFile, defaultConfigFile)
		d.log.Debug("Deployment file not specified, using {}", defaultConfigFile)
	}

	_, err := os.Stat(configFile)
	if err != nil && !os.IsExist(err) {
		d.log.Fatal("No deployment config file exists at: {}", configFile)
	}

	contentstring, err := ioutil.ReadFile(configFile)
	if err != nil {
		d.log.Fatal("Deployment config file could not be read: {}", err)
	}

	if ok, err := utils.IsYaml(contentstring); !ok {
		d.log.Fatal("Deployment config file ({}) is not valid YAML: {}", configFile, err)
	}

	err = yaml.Unmarshal([]byte(contentstring), &d.config)
	if err != nil {
		d.log.Fatal("Error parsing deployment config {}", err)
	}

	d.config.configFilePath = configFile

	d.processConfig()

}

// processConfig ensures that the deployment config is valid
func (d *Deploy) processConfig() {

	// Set defaults
	setConfigDefault(&d.config.Deployment.Container.Repo, defaultContainerRepo)
	setConfigDefault(&d.config.Deployment.Container.Tag, defaultContainerTag)
	setConfigDefault(&d.config.Deployment.Directory, defaultDeployDirectory)
	setConfigDefault(&d.config.Deployment.Script, defaultDeployScript)

	// Create our global spec if it doesn't exist so we don't have to keep checking if it exists
	if d.config.Global.Spec == nil {
		d.config.Global.Spec = &Spec{}
	}

	d.validateSpec(d.config.Global.Spec)

	d.config.environmentMap = make(map[string]int)
	for i, environment := range d.config.Environments {

		// Check to make sure that we don't have multiple environments with the same name
		if _, ok := d.config.environmentMap[environment.Name]; ok {
			d.log.Fatal("Error parsing config, duplicate environment name `{}` found", environment.Name)
		}

		// Ensure there are instances for this environment
		if len(environment.Instances) <= 0 {
			d.log.Fatal("No instances found for environment: `{}`", environment.Name)
		}

		d.config.environmentMap[environment.Name] = i

		// Create our environment spec if it doesn't exist so we don't have to keep checking if it exists
		if environment.Spec == nil {
			environment.Spec = &Spec{}
		}

		d.validateSpec(environment.Spec)

		environment.instanceMap = make(map[string]int)
		for j, instance := range environment.Instances {

			// Check to make sure that we don't have multiple instances with the same name
			if _, ok := environment.instanceMap[instance.Name]; ok {
				d.log.Fatal("Error parsing config, duplicate instance name '{}' for environment '{}'", instance.Name, environment.Name)
			}

			// Ensure the instance name does not conflict with the ALL option name.  This is a reserved name for designating a deployment to all instances in an environment via the manual prompt list
			if strings.ToLower(instance.Name) == strings.ToLower(allOptionPrompt) || strings.ToLower(instance.Name) == strings.ToLower(allOptionCli) {
				d.log.Fatal("Deployment config cannot have an instance named '{}'. It is a reserved name.", instance.Name)
			}

			environment.instanceMap[instance.Name] = j

			// Create our instance spec if it doesn't exist so we don't have to keep checking if it exists
			if instance.Spec == nil {
				instance.Spec = &Spec{}
			}

			d.validateSpec(instance.Spec)

			// Merge all of the secrets and environment variables
			// Instance-level specs take precedence, followed by environment-level then global-level
			if instance.Spec.Kubernetes.ServiceAccount == "" {
				if environment.Spec.Kubernetes.ServiceAccount != "" {
					instance.Spec.Kubernetes.ServiceAccount = environment.Spec.Kubernetes.ServiceAccount
				} else if d.config.Global.Spec.Kubernetes.ServiceAccount != "" {
					instance.Spec.Kubernetes.ServiceAccount = d.config.Global.Spec.Kubernetes.ServiceAccount
				} else {
					d.log.Fatal("Kubernetes service account is not set for instance '{}' in environment '{}'", instance.Name, environment.Name)
				}
			}
			if instance.Spec.Kubernetes.Cluster == "" {
				if environment.Spec.Kubernetes.Cluster != "" {
					instance.Spec.Kubernetes.Cluster = environment.Spec.Kubernetes.Cluster
				} else if d.config.Global.Spec.Kubernetes.Cluster != "" {
					instance.Spec.Kubernetes.Cluster = d.config.Global.Spec.Kubernetes.Cluster
				} else {
					d.log.Fatal("Kubernetes cluster is not set for instance '{}' in environment '{}'", instance.Name, environment.Name)
				}
			}

			instance.Spec.Tools = mergeTools(instance.Spec.Tools, environment.Spec.Tools, d.config.Global.Spec.Tools)
			instance.Spec.EnvironmentVars = mergeEnvVars(instance.Spec.EnvironmentVars, environment.Spec.EnvironmentVars, d.config.Global.Spec.EnvironmentVars)
                        instance.Spec.EnvironmentVars = helmifyDoSets(instance.Spec.EnvironmentVars)
			instance.Spec.Secrets = mergeSecrets(instance.Spec.Secrets, environment.Spec.Secrets, d.config.Global.Spec.Secrets)

			// Get Vault details
			vault := d.stim.Vault()
			vaultToken, err := vault.GetToken()
			if err != nil {
				d.log.Fatal("Error fetching Vault token for deploy '{}'", err)
			}

			vaultAddress, err := vault.GetAddress()
			if err != nil {
				d.log.Fatal("Error fetching Vault address for deploy '{}'", err)
			}

			// Generate stim env vars
			stimEnvs := []*EnvironmentVar{}

			stimEnvs = append(stimEnvs, []*EnvironmentVar{
				&EnvironmentVar{Name: "VAULT_ADDR", Value: vaultAddress},
				&EnvironmentVar{Name: "VAULT_TOKEN", Value: vaultToken},
				&EnvironmentVar{Name: "DEPLOY_ENVIRONMENT", Value: environment.Name},
				&EnvironmentVar{Name: "DEPLOY_INSTANCE", Value: instance.Name},
				&EnvironmentVar{Name: "DEPLOY_CLUSTER", Value: instance.Spec.Kubernetes.Cluster},
			}...)

			// Generate the Kube config secret
			var stimSecrets []*v2e.SecretItem
			secretMap := make(map[string]string)
			secretMap["CLUSTER_SERVER"] = "cluster-server"
			secretMap["CLUSTER_CA"] = "cluster-ca"
			secretMap["USER_TOKEN"] = "user-token"
			stimSecrets = append(stimSecrets, &v2e.SecretItem{
				SecretPath: fmt.Sprintf("secret/kubernetes/%s/%s/kube-config", instance.Spec.Kubernetes.Cluster, instance.Spec.Kubernetes.ServiceAccount),
				SecretMaps: secretMap,
			})

			// Add stim envs/secrets and ensure no reserved env vars have been set
			d.finalizeEnv(instance, stimEnvs, stimSecrets)
		}
	}

	// Determine the full directory path
	configAbs, err := filepath.Abs(d.config.configFilePath)
	if err != nil {
		d.log.Fatal("Error fetching deploy filepath '{}'", err)
	}
	d.config.Deployment.fullDirectoryPath = filepath.Join(filepath.Dir(configAbs), d.config.Deployment.Directory)
}

// Generate the list of reserved env var names
func (d *Deploy) finalizeEnv(instance *Instance, stimEnvs []*EnvironmentVar, stimSecrets []*v2e.SecretItem) {

	// Generate the list of reserved env var names (additionally SECRET_CONFIG as we'll add that one at the end)
	reservedVarNames := []string{"SECRET_CONFIG", "STIM_DEPLOY"}

	for _, s := range stimEnvs {
		reservedVarNames = append(reservedVarNames, s.Name)
	}
	for _, s := range stimSecrets {
		for m := range s.SecretMaps {
			reservedVarNames = append(reservedVarNames, m)
		}
	}

	// Exit if any user-provided environment vars conflict with reserved ones
	for _, e := range instance.Spec.EnvironmentVars {
		if utils.Contains(reservedVarNames, e.Name) {
			d.log.Fatal("Reserved environment variable name '{}' found in config", e.Name)
		}
	}
	for _, s := range instance.Spec.Secrets {
		for m := range s.SecretMaps {
			if utils.Contains(reservedVarNames, m) {
				d.log.Fatal("Reserved environment variable name '{}' found in config", m)
			}
		}
	}

	// Combine our secrets
	instance.Spec.Secrets = append(instance.Spec.Secrets, stimSecrets...)

	// Create the secret config
	secretConfig, err := d.makeSecretConfig(instance)
	if err != nil {
		d.log.Fatal("Error making secret config '{}'", err)
	}
	stimEnvs = append(stimEnvs, &EnvironmentVar{Name: "SECRET_CONFIG", Value: secretConfig})
	stimEnvs = append(stimEnvs, &EnvironmentVar{Name: "STIM_DEPLOY", Value: "true"})

	// Combine our env vars
	instance.Spec.EnvironmentVars = append(instance.Spec.EnvironmentVars, stimEnvs...)

}

// validateSpec validates fields in a config 'spec' section to ensure that it
// meets all requirements
func (d *Deploy) validateSpec(spec *Spec) {
	for toolName, toolSpec := range spec.Tools {
		if toolName == "helm" && toolSpec.Version == "" {
			d.log.Fatal("Version detection not supported for helm, please specify a version in the `spec.tools.helm` config")
		}
	}
}

// helmifyDoSets looks for any variables starting with "STIM_HELM" and converts them into giant '--set' list and returns string.
func helmifyDoSets(instance []*EnvironmentVar) []*EnvironmentVar {

      var slug = defaultHelmifyPrefix
      envMap := make(map[string]string)
      setSlice := []string{}
      result := instance

      for _, s := range instance {
            if strings.Contains(s.Name, slug) {
                  if _, ok := envMap[s.Name]; !ok {
                        envMap[s.Name] = s.Value
                  }
            }
      }
      for k, v := range envMap {
            var command = strings.TrimPrefix(k, slug)
            command = fmt.Sprintf("--set %s=%s", command, v)
            setSlice = append(setSlice, command)
      }
      v := new(EnvironmentVar)
      v.Name = defaultHelmifySlug
      v.Value = strings.Join(setSlice, " \\ \n")
      result = append(result, v)
      return result
}

// mergeEnvVars is used to merge environment variable configuration at the various levels it can be set at
func mergeEnvVars(instance []*EnvironmentVar, environment []*EnvironmentVar, global []*EnvironmentVar) []*EnvironmentVar {

	result := instance

	// Add environment envVars (if they don't already exist)
	for _, e := range environment {
		exists := false
		for _, inst := range result {
			if inst.Name == e.Name {
				exists = true
			}
		}

		// Add the item if it doesn't exist
		if !exists {
			result = append(result, e)
		}
	}

	// Add global envVars (if they don't already exist)
	for _, g := range global {
		exists := false
		for _, inst := range result {
			if inst.Name == g.Name {
				exists = true
			}
		}

		// Add the item if it doesn't exist
		if !exists {
			result = append(result, g)
		}
	}

	return result
}

// mergeSecrets is used to merge secret configs at the various levels they can be set at
func mergeSecrets(instance []*v2e.SecretItem, environment []*v2e.SecretItem, global []*v2e.SecretItem) []*v2e.SecretItem {

	result := global

	// Add environment envVars
	for _, e := range environment {
		result = append(result, e)
	}

	// Add global envVars to instance (if they don't already exist)
	for _, inst := range instance {
		result = append(result, inst)
	}

	return result
}

// mergeTools is used to merge tool configurations
func mergeTools(instance map[string]stim.EnvTool, environment map[string]stim.EnvTool, global map[string]stim.EnvTool) map[string]stim.EnvTool {

	result := make(map[string]stim.EnvTool)

	// Set Global tools
	for k, v := range global {
		result[k] = v
	}

	// Overwrite with instance tools
	for k, v := range environment {
		if v.Unset == true {
			delete(result, k)
		} else {
			result[k] = v
		}
	}

	// Overwrite with instance tools
	for k, v := range instance {
		if v.Unset == true {
			delete(result, k)
		} else {
			result[k] = v
		}
	}

	return result
}

// setConfigDefault is used to set a default value (if it doesn't exist)
func setConfigDefault(value *string, def string) {
	if len(*value) == 0 {
		*value = def
	}
}

package deploy

import (
	"strings"

	log "github.com/PremiereGlobal/stim/pkg/stimlog"
	"github.com/PremiereGlobal/stim/stim"
)

const (
	allOptionPrompt = "--ALL--"
	allOptionCli    = "all"
)

// Deploy is the primary type for the stim deploy subcommand
type Deploy struct {
	name   string
	stim   *stim.Stim
	config Config
	log    log.StimLogger
}

// New creates a new 'Deploy' object
func New() *Deploy {
	deploy := &Deploy{}
	return deploy
}

// Name is a required stim function that returns the name of the stimpack
func (d *Deploy) Name() string {
	return d.name
}

// Run is the main entrypoint to the "deploy" command
func (d *Deploy) Run() {

	d.log = d.stim.GetLogger()

	// Read in the config file and set up defaults
	d.parseConfig()

	// Determine the selected environment (via cli param) or prompt the user
	selectedEnvironmentName := ""
	environmentArg := d.stim.ConfigGetString("deploy.environment")
	if environmentArg != "" {
		if _, ok := d.config.environmentMap[environmentArg]; ok {
			selectedEnvironmentName = environmentArg
		} else {
			d.log.Fatal("Provided environment value '{}' is not in config file", environmentArg)
		}
	} else {
		environmentList := make([]string, len(d.config.Environments))
		for i, e := range d.config.Environments {
			environmentList[i] = e.Name
		}
		selectedEnvironmentName, _ = d.stim.PromptList("Which environment?", environmentList, d.stim.ConfigGetString("deploy.environment"))
	}
	selectedEnvironment := d.config.Environments[d.config.environmentMap[selectedEnvironmentName]]

	// Determine the selected instance (via cli param) or prompt the user
	instanceList := make([]string, len(selectedEnvironment.Instances)+1)
	instanceList[0] = allOptionPrompt
	for i, inst := range selectedEnvironment.Instances {
		instanceList[i+1] = inst.Name
	}
	selectedInstanceName, _ := d.stim.PromptList("Which instance?", instanceList, d.stim.ConfigGetString("deploy.instance"))
	if strings.ToLower(selectedInstanceName) == strings.ToLower(allOptionPrompt) || strings.ToLower(selectedInstanceName) == strings.ToLower(allOptionCli) {
		selectedInstanceName = allOptionCli
	} else if _, ok := selectedEnvironment.instanceMap[selectedInstanceName]; !ok {
		d.log.Fatal("Provided instance value '{}' is not in config file under environment '{}'", selectedInstanceName, selectedEnvironmentName)
	}

	// Run the deployment(s)
	if selectedInstanceName == allOptionCli {
		d.log.Info("Deploying to all clusters in environment: {}", selectedEnvironment.Name)
		for _, inst := range selectedEnvironment.Instances {
			d.Deploy(selectedEnvironment, inst)
		}
	} else {
		d.Deploy(selectedEnvironment, selectedEnvironment.Instances[selectedEnvironment.instanceMap[selectedInstanceName]])
	}

}

// Deploy runs the deployment in the way that the user wants
func (d *Deploy) Deploy(environment *Environment, instance *Instance) {

	d.log.Info("Deploying to '{}' environment in instance: {}", environment.Name, instance.Name)

	// For now, only the kube-vault-deploy docker method is implemented but more could be added here...
	d.startDeployContainer(instance)

}

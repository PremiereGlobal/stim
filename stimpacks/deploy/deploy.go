package deploy

import (
	"errors"
	"fmt"

	"github.com/PremiereGlobal/stim/stim"
	// "github.com/fatih/color"
)

const (
	ALL_OPTION_TEXT = "--ALL--"
)

type Deploy struct {
	name   string
	stim   *stim.Stim
	config Config
}

func New() *Deploy {
	deploy := &Deploy{}
	return deploy
}

func (d *Deploy) Name() string {
	return d.name
}

// Run is the main entrypoint to the "deploy" command
func (d *Deploy) Run() {

	// Read in the config file and set up defaults
	d.parseConfig()

	// Determine the selected environment (via cli param) or prompt the user
	selectedEnvironmentName := ""
	environmentArg := d.stim.GetConfig("deploy.environment")
	if environmentArg != "" {
		if _, ok := d.config.environmentMap[environmentArg]; ok {
			selectedEnvironmentName = environmentArg
		} else {
			d.stim.Fatal(errors.New(fmt.Sprintf("Provided environment value '%s' is not in config file", environmentArg)))
		}
	} else {
		environmentList := make([]string, len(d.config.Environments))
		for i, e := range d.config.Environments {
			environmentList[i] = e.Name
		}
		selectedEnvironmentName, _ = d.stim.PromptList("Which environment?", environmentList, d.stim.GetConfig("deploy.environment"))
	}
	selectedEnvironment := d.config.Environments[d.config.environmentMap[selectedEnvironmentName]]

	// Determine the selected instance (via cli param) or prompt the user
	selectedInstanceName := ""
	instanceArg := d.stim.GetConfig("deploy.instance")
	if instanceArg != "" {
		if _, ok := selectedEnvironment.instanceMap[instanceArg]; ok {
			selectedInstanceName = instanceArg
		} else {
			d.stim.Fatal(errors.New(fmt.Sprintf("Provided instance value '%s' is not in config file under environment '%s'", instanceArg, selectedEnvironmentName)))
		}
	} else {
		instanceList := make([]string, len(selectedEnvironment.Instances)+1)
		instanceList[0] = ALL_OPTION_TEXT
		for i, inst := range selectedEnvironment.Instances {
			instanceList[i+1] = inst.Name
		}

		selectedInstanceName, _ = d.stim.PromptList("Which instance?", instanceList, d.stim.GetConfig("deploy.instance"))
	}

	// Run the deployment(s)
	// color.Set(color.FgGreen)
	if selectedInstanceName == ALL_OPTION_TEXT {
		fmt.Println(fmt.Sprintf("Deploying to all clusters in environment: %s", selectedEnvironment.Name))
		for _, inst := range selectedEnvironment.Instances {
			d.Deploy(selectedEnvironment, inst)
		}
	} else {
		d.Deploy(selectedEnvironment, selectedEnvironment.Instances[selectedEnvironment.instanceMap[selectedInstanceName]])
	}

}

// Run the deployment in the way that the user wants
func (d *Deploy) Deploy(environment *Environment, instance *Instance) {

	fmt.Println(fmt.Sprintf("Deploying to '%s' environment in instance: %s", environment.Name, instance.Name))

	// For now, only the kube-vault-deploy docker method is implemented but more could be added here...
	d.startDeployContainer(environment, instance)

}

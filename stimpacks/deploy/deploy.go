package deploy

import (
	"os"
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
		if selectedEnvironmentName == "" {
			d.log.Info("No environment selected! exiting")
			os.Exit(0)
		}
	}
	selectedEnvironment := d.config.Environments[d.config.environmentMap[selectedEnvironmentName]]

	// Determine the selected instance (via cli param) or prompt the user
	instanceList := make([]string, 0)

	//Check if we should remove all prompt or not
	if !selectedEnvironment.RemoveAllPrompt {
		instanceList = append(instanceList, allOptionPrompt)
	}
	for _, inst := range selectedEnvironment.Instances {
		instanceList = append(instanceList, inst.Name)
	}
	selectedInstanceName, _ := d.stim.PromptList("Which instance?", instanceList, d.stim.ConfigGetString("deploy.instance"))
	if selectedInstanceName == "" {
		d.log.Info("No instance selected! exiting")
		os.Exit(0)
	}
	if strings.ToLower(selectedInstanceName) == strings.ToLower(allOptionPrompt) || strings.ToLower(selectedInstanceName) == strings.ToLower(allOptionCli) {
		selectedInstanceName = allOptionCli
	} else if _, ok := selectedEnvironment.instanceMap[selectedInstanceName]; !ok {
		d.log.Fatal("Provided instance value '{}' is not in config file under environment '{}'", selectedInstanceName, selectedEnvironmentName)
	}

	// Run the deployment(s)
	if selectedInstanceName == allOptionCli {
		d.log.Info("Deploying to all clusters in environment: {}", selectedEnvironment.Name)
		//Check if confirmation prompt is required
		if selectedEnvironment.Spec.AddConfirmationPrompt {
			//Do AddConfirmationPrompt, only if the instance is not passed on the cli
			proceed, _ := d.stim.PromptBool("Proceed?", d.stim.ConfigGetString("deploy.instance") != "", false)
			if !proceed {
				os.Exit(1)
			}
		}
		for _, inst := range selectedEnvironment.Instances {
			if inst.Spec.AddConfirmationPrompt {
				//Do AddConfirmationPrompt, only if the instance is not passed on the cli
				proceed, _ := d.stim.PromptBool("Proceed?", d.stim.ConfigGetString("deploy.instance") != "", false)
				if !proceed {
					os.Exit(1)
				}
			}
			d.Deploy(selectedEnvironment, inst)
		}
	} else {
		d.log.Info("Deploying to environment: {} and instance: {}", selectedInstanceName)
		inst := selectedEnvironment.Instances[selectedEnvironment.instanceMap[selectedInstanceName]]
		if selectedEnvironment.Spec.AddConfirmationPrompt || inst.Spec.AddConfirmationPrompt {
			proceed, _ := d.stim.PromptBool("Proceed?", d.stim.ConfigGetString("deploy.instance") != "", false)
			if !proceed {
				os.Exit(1)
			}
		}
		d.Deploy(selectedEnvironment, inst)
	}

}

// Deploy runs the deployment in the way that the user wants
func (d *Deploy) Deploy(environment *Environment, instance *Instance) {

	d.log.Info("Deploying to '{}' environment in instance: {}", environment.Name, instance.Name)

	// For now, only the kube-vault-deploy docker method is implemented but more could be added here...
	d.startDeployContainer(instance)

}

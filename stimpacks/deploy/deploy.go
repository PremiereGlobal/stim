package deploy

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/PremiereGlobal/stim/pkg/docker"
	log "github.com/PremiereGlobal/stim/pkg/stimlog"
	"github.com/PremiereGlobal/stim/stim"
	"golang.org/x/mod/semver"
)

const (
	allOptionPrompt = "--ALL--"
	allOptionCli    = "all"
)

const (
	DEPLOY_METHOD_UNKNOWN int = 0
	DEPLOY_METHOD_DOCKER  int = 1
	DEPLOY_METHOD_SHELL   int = 2
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
	if d.config.Global.RequiredVersion != "" {
		d.log.Info("Deploy has set RequiredVersion to:{}, currently:{}", d.config.Global.RequiredVersion, d.stim.GetVersion())
		if semver.Compare(d.stim.GetVersion(), d.config.Global.RequiredVersion) != 0 {
			d.log.Fatal("Stim is not at the Required version for deploy, current:{}, required:{}\n\t Please check https://github.com/PremiereGlobal/stim/releases for new versions", d.stim.GetVersion(), d.config.Global.RequiredVersion)
		}
	} else {
		if d.config.Global.MinimumVersion != "" {
			d.log.Info("Deploy has set MinimumVersion to:{}, currently:{}", d.config.Global.MinimumVersion, d.stim.GetVersion())
			if semver.Compare(d.stim.GetVersion(), d.config.Global.MinimumVersion) < 0 {
				d.log.Fatal("Stim is not at the Required version for deploy, current:{}, minimum:{}\n\t Please check https://github.com/PremiereGlobal/stim/releases for new versions", d.stim.GetVersion(), d.config.Global.MinimumVersion)
			}
		}
	}

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

	deployMethod, err := d.DetermineDeployMethod()
	if err != nil {
		d.log.Fatal(err)
	}

	if deployMethod == DEPLOY_METHOD_DOCKER {
		d.startDeployContainer(instance)
	} else if deployMethod == DEPLOY_METHOD_SHELL {
		d.startDeployShell(instance)
	} else {
		d.log.Fatal("Could not determine deployment method")
	}

}

// DetermineDeployMethod figures out the deploy method based on user input
// and availability
func (d *Deploy) DetermineDeployMethod() (int, error) {

	deployMethod := d.stim.ConfigGetString("deploy.method")

	isInDocker := docker.IsInDocker()
	isDockerAvailable, _ := docker.IsDockerAvailable()

	if deployMethod == "auto" && isDockerAvailable && !isInDocker {
		d.log.Debug("Using Docker to deploy (auto)")
		return DEPLOY_METHOD_DOCKER, nil
	}

	if deployMethod == "docker" && isDockerAvailable && !isInDocker {
		d.log.Debug("Using Docker to deploy (specified by user)")
		return DEPLOY_METHOD_DOCKER, nil
	}

	// If we are already in a container, only shell is supported
	if deployMethod == "auto" && isInDocker {
		d.log.Debug("Using shell to deploy (auto) as detected we're running in Docker")
		return DEPLOY_METHOD_SHELL, nil
	}

	// If docker is not available, and auto is selected, force the user to specify --shell
	// This is to avoid inadvertently running shell commands on their machine
	if deployMethod == "auto" && !isDockerAvailable {
		return DEPLOY_METHOD_UNKNOWN, errors.New("Docker is not available.  To deploy using shell use '--method=shell' argument (this is not recommended)")
	}

	if deployMethod == "shell" {
		d.log.Debug("Using shell to deploy (specified by user)")
		return DEPLOY_METHOD_SHELL, nil
	}

	// Below we're detecting some specific error cases to give more info to the user

	if deployMethod == "docker" && isInDocker {
		return DEPLOY_METHOD_UNKNOWN, errors.New("Cannot deploy with Docker as we are already in a container")
	}
	if deployMethod == "docker" && !isDockerAvailable {
		return DEPLOY_METHOD_UNKNOWN, errors.New("Cannot deploy with Docker as it is not available")
	}

	return DEPLOY_METHOD_UNKNOWN, errors.New(fmt.Sprintf("Invalid deployment method '%s' provided.  Must be one of ['auto','docker','shell']", deployMethod))
}

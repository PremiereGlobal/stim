package deploy

import (
	"fmt"

	"github.com/PremiereGlobal/stim/stim"
)

// startDeployShell starts an instance deployment using the command shell
func (d *Deploy) startDeployShell(instance *Instance) {

	envs := make([]string, len(instance.Spec.EnvironmentVars))
	for i, e := range instance.Spec.EnvironmentVars {
		envs[i] = fmt.Sprintf("%s=%s", e.Name, e.Value)
	}

	d.log.Debug("Setting working directory {}", d.config.Deployment.fullDirectoryPath)
	e := d.stim.Env(&stim.EnvConfig{
		EnvVars: envs,
		Kubernetes: &stim.EnvConfigKubernetes{
			Cluster:          instance.Spec.Kubernetes.Cluster,
			ServiceAccount:   instance.Spec.Kubernetes.ServiceAccount,
			DefaultNamespace: "default"},
		Vault: &stim.EnvConfigVault{
			SecretItems: instance.Spec.Secrets,
		},
		WorkDir: d.config.Deployment.fullDirectoryPath,
		Tools:   instance.Spec.Tools,
	})

	d.log.Debug("Running script ./{}", d.config.Deployment.Script)
	out, err := e.Run("./" + d.config.Deployment.Script)
	if err != nil {
		d.log.Fatal("Error running command: {}", err)
	}

	d.log.Info(out)
}

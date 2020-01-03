package deploy

import (
	"fmt"
	// "path/filepath"

	"github.com/PremiereGlobal/stim/stim"
)

func (d *Deploy) startDeployShell(instance *Instance) {

	envs := make([]string, len(instance.Spec.EnvironmentVars))
	for i, e := range instance.Spec.EnvironmentVars {
		envs[i] = fmt.Sprintf("%s=%s", e.Name, e.Value)
	}

	// e := stim.NewEnv()
	// e.SetEnvs(envs)
	// e.SetSecrets(instance.Spec.Secrets)
	// e.SetKubeCluster(instance.Spec.Kubernetes.Cluster, instance.Spec.Kubernetes.ServiceAccount)

	e := d.stim.Env(&stim.EnvConfig{
		EnvVars: envs,
		Kubernetes: &stim.EnvConfigKubernetes{
			Cluster:          instance.Spec.Kubernetes.Cluster,
			ServiceAccount:   instance.Spec.Kubernetes.ServiceAccount,
			DefaultNamespace: "kube-system"},
		Vault: &stim.EnvConfigVault{
			SecretItems: instance.Spec.Secrets,
		},
		WorkDir: d.config.Deployment.fullDirectoryPath,
	})

	d.log.Debug("Setting working directory {}", d.config.Deployment.fullDirectoryPath)
	// err := e.SetWorkDir()
	// if err != nil {
	// 	d.log.Fatal("Set work dir error {}", err)
	// }

	// d.log.Info(e.Run("cat $KUBECONFIG"))
	d.log.Debug("Running script ./{}", d.config.Deployment.Script)
	// d.log.Info(e.Run("ls "+filepath.Join(d.config.Deployment.fullDirectoryPath, d.config.Deployment.Script)))

	// a,_ := e.Run("ls -al; pwd; env")
	// d.log.Info(a)

	// out, err := e.Run("./"+d.config.Deployment.Script)
	out, err := e.Run("./" + d.config.Deployment.Script)
	if err != nil {
		d.log.Fatal("Error running command: {}", err)
	}

	d.log.Info(out)
}

//
// func (d *Deploy) getBinVersionKube() string {
//
//   // See if the binary is even available
//   if !d.execBinAvailable("kubectl") {
//     d.log.Debug("Executable kubectl not found")
//     return ""
//   }
//
//   kubeVersionFull := strings.Trim(d.runShellCmd("kubectl version --client --short"), "\n")
//   r, err := regexp.Compile("\\: (.*)")
//   if err != nil {
//     d.log.Fatal("Unable to compile regex for Kubernetes binary version detection '{}'", err)
//   }
//   kubeVersionParts := r.FindStringSubmatch(kubeVersionFull)
//
//   if len(kubeVersionParts) != 2 {
//     d.log.Fatal("Unable to determine Kubernetes binary version. Output of command is '{}'", kubeVersionFull)
//   }
//
//   return kubeVersionParts[1]
// }
//
// func (d *Deploy) getBinVersionHelm() string {
//
//   // See if the binary is even available
//   if !d.execBinAvailable("helm") {
//     d.log.Debug("Executable helm not found")
//     return ""
//   }
//
//   helmVersion := d.runShellCmd("helm version --client --template {{.Client.SemVer}}")
//
//   return helmVersion
// }
//
// func (d *Deploy) getBinVersionKops() string {
//
//   // See if the binary is even available
//   if !d.execBinAvailable("kops") {
//     d.log.Debug("Executable kops not found")
//     return ""
//   }
//
//   kopsVersionFull := d.runShellCmd("kops version")
//   r, err := regexp.Compile("Version (.*) .*")
//   if err != nil {
//     d.log.Fatal("Unable to compile regex for Kops binary version detection '{}'", err)
//   }
//   kopsVersionParts := r.FindStringSubmatch(kopsVersionFull)
//
//   if len(kopsVersionParts) != 2 {
//     d.log.Fatal("Unable to determine Kops binary version. Output of command is '{}'", kopsVersionFull)
//   }
//
//   return kopsVersionParts[1]
// }
//
// func (d *Deploy) getBinVersionVault() string {
//
//   // See if the binary is even available
//   if !d.execBinAvailable("vault") {
//     d.log.Debug("Executable vault not found")
//     return ""
//   }
//
//   vaultVersionFull := strings.Trim(d.runShellCmd("vault version"), "\n")
//   r, err := regexp.Compile("Vault (.*) .*")
//   if err != nil {
//     d.log.Fatal("Unable to compile regex for Vault binary version detection '{}'", err)
//   }
//   vaultVersionParts := r.FindStringSubmatch(vaultVersionFull)
//
//   if len(vaultVersionParts) != 2 {
//     d.log.Fatal("Unable to determine Vault binary version. Output of command is '{}'", vaultVersionFull)
//   }
//
//   return vaultVersionParts[1]
// }

package stim

import (
	"fmt"
	"path/filepath"

	"github.com/PremiereGlobal/stim/pkg/env"
	"github.com/PremiereGlobal/stim/pkg/kubernetes"
	"github.com/PremiereGlobal/vault-to-envs/pkg/vaulttoenvs"
)

type EnvConfig struct {
	EnvVars    []string
	Kubernetes *EnvConfigKubernetes
	Vault      *EnvConfigVault
	WorkDir    string
}

type EnvConfigKubernetes struct {
	Cluster          string
	ServiceAccount   string
	DefaultNamespace string
}

type EnvConfigVault struct {
	SecretItems     []*vaulttoenvs.SecretItem
	CliVersionVault string
}

func (stim *Stim) Env(config *EnvConfig) *env.Env {

	e, err := env.New(env.Config{})
	if err != nil {
		stim.log.Fatal("Stim: Error creating new environment.", err)
	}

	e.SetWorkDir(config.WorkDir)

	e.AddEnvVars(config.EnvVars...)

	// If requiring Kubernetes, set things up
	if config.Kubernetes != nil {

		// This is the path where the kubeconfig will be written
		kubeConfigFilePath := filepath.Join(e.GetPath(), "kubeconfig")

		vault := stim.Vault()

		// Get the Kubernetes creds from Vault
		secretValues, err := vault.GetSecretKeys("secret/kubernetes/" + config.Kubernetes.Cluster + "/" + config.Kubernetes.ServiceAccount + "/kube-config")
		if err != nil {
			stim.log.Fatal("Stim: Error getting kubeconfig secrets for environment. {}", err)
		}

		// If namespace not set use the default from Vault
		defaultNamespace := config.Kubernetes.DefaultNamespace
		if defaultNamespace == "" {
			defaultNamespace = secretValues["default-namespace"]
		}

		// Build the Kube config options
		kubeConfigOptions := &kubernetes.KubeConfigOptions{
			ClusterName:             config.Kubernetes.Cluster,
			ClusterServer:           secretValues["cluster-server"],
			ClusterCA:               secretValues["cluster-ca"],
			AuthName:                config.Kubernetes.Cluster + "-" + config.Kubernetes.ServiceAccount,
			AuthToken:               secretValues["user-token"],
			ContextName:             config.Kubernetes.Cluster,
			ContextSetCurrent:       true,
			ContextDefaultNamespace: defaultNamespace,
			// KubeConfigFilePath:      kubeConfigFilePath,
		}

		kc := kubernetes.NewKubeConfig(kubeConfigFilePath)
		// kc.SetKubeconfig(kubeConfigOptions)

		// Write out the kubeconfig file
		// kube := stim.Kubernetes()
		err = kc.ModifyConfig(kubeConfigOptions)
		if err != nil {
			stim.log.Fatal("Stim: Error creating kubeconfig for environment. {}", err)
		}

		// Tell the environment to use the kubeconfig in the environment PATH
		e.AddEnvVars([]string{fmt.Sprintf("%s=%s", "KUBECONFIG", kubeConfigFilePath)}...)

		// Get the verison of Kubernetes
		// kube.

		// Get the version of Helm
		// kubectl get po -n ${TILLER_NAMESPACE} -l app=helm,name=tiller

		// Link to our Kubernetes version
		// e.PathLink("cach-dir/kubectl-v1.10.8", "kubectl")
	}

	// If requiring secrets, set those up
	if config.Vault != nil && len(config.Vault.SecretItems) > 0 {

		vault := stim.Vault()

		vaultAddress, err := vault.GetAddress()
		if err != nil {
			stim.log.Fatal("Stim: Unable to get Vault address for environment. {}", err)
		}

		vaultToken, err := vault.GetToken()
		if err != nil {
			stim.log.Fatal("Stim: Unable to get Vault token for environment. {}", err)
		}

		v2e := vaulttoenvs.NewVaultToEnvs(&vaulttoenvs.Config{
			VaultAddr: vaultAddress,
		})
		v2e.SetVaultToken(vaultToken)
		v2e.AddSecretItems(config.Vault.SecretItems...)

		secretEnvs, err := v2e.GetEnvs()
		if err != nil {
			stim.log.Fatal("Stim: Unable to get Vault secrets for environment. {}", err)
		}

		e.AddEnvVars(secretEnvs...)

	}

	return e
}

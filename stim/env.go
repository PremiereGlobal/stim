package stim

import (
	"fmt"
	"path/filepath"

	"github.com/PremiereGlobal/stim/pkg/env"
	"github.com/PremiereGlobal/stim/pkg/kubernetes"
	v2e "github.com/PremiereGlobal/vault-to-envs/pkg/vaulttoenvs"
)

type EnvConfig struct {
	EnvVars    []string
	Kubernetes *EnvConfigKubernetes
	Vault      *EnvConfigVault
}

type EnvConfigKubernetes struct {
	Cluster          string
	ServiceAccount   string
	DefaultNamespace string
}

type EnvConfigVault struct {
	SecretConfig    []*v2e.SecretItem
	CliVersionVault string
}

func (stim *Stim) Env(config *EnvConfig) *env.Env {

	e, err := env.New(env.Config{})
	if err != nil {
		stim.log.Fatal("Stim: Error creating new environment.", err)
	}

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
		if defaultNamespace != "" {
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
			KubeConfigFilePath:      kubeConfigFilePath,
		}

		// Write out the kubeconfig file
		kube := stim.Kubernetes()
		err = kube.SetKubeconfig(kubeConfigOptions)
		if err != nil {
			stim.log.Fatal("Stim: Error creating kubeconfig for environment. {}", err)
		}

		// Tell the environment to use the kubeconfig in the environment PATH
		e.AddEnvVars([]string{fmt.Sprintf("%s=%s", "KUBECONFIG", kubeConfigFilePath)}...)

		// Link to our Kubernetes version
		// e.PathLink("cach-dir/kubectl-v1.10.8", "kubectl")
	}

	return e
}

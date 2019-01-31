package kubernetes

import (
	kubepkg "github.com/readytalk/stim/pkg/kubernetes"
	// "github.com/davecgh/go-spew/spew"
)

func (k *Kubernetes) configureContext() error {

	// Create a Vault instance
	k.vault = k.stim.Vault()

	cluster, err := k.stim.PromptListVault("secret/kubernetes", "Select Cluster", k.stim.GetConfig("kube-config-cluster"))
	if err != nil {
		return err
	}

	sa, err := k.stim.PromptListVault("secret/kubernetes/"+cluster, "Select Service Account", k.stim.GetConfig("kube-service-account"))
	if err != nil {
		return err
	}

	// Get secrets from Vault
	secretValues, err := k.vault.GetSecretKeys("secret/kubernetes/" + cluster + "/" + sa + "/kube-config")
	if err != nil {
		return err
	}

	namespace, err := k.stim.PromptString("Select Default Namespace", k.stim.GetConfig("kube-service-account"), secretValues["default-namespace"])
	if err != nil {
		return err
	}

	context, err := k.stim.PromptString("Context Name", k.stim.GetConfig("kube-context"), cluster)
	if err != nil {
		return err
	}

	currentContext, err := k.stim.PromptBool("Set as current context?", k.stim.GetConfigBool("kube-current-context"), true)
	if err != nil {
		return err
	}

	// Build the config options
	kubeConfigOptions := &kubepkg.KubeConfigOptions{
		ClusterName:             cluster,
		ClusterServer:           secretValues["cluster-server"],
		ClusterCA:               secretValues["cluster-ca"],
		AuthName:                cluster + "-" + sa,
		AuthToken:               secretValues["user-token"],
		ContextName:             context,
		ContextSetCurrent:       currentContext,
		ContextDefaultNamespace: namespace,
	}

	// Set the kubeconfig
	kube := k.stim.Kubernetes()
	err = kube.SetKubeconfig(kubeConfigOptions)
	if err != nil {
		return err
	}

	return nil
}

package kubernetes

import (
	"errors"
	kubepkg "github.com/readytalk/stim/pkg/kubernetes"
)

func (k *Kubernetes) configureContext() error {

	// Ensure we have the fields we need
	err := k.validateConfigureFields()
	if err != nil {
		return err
	}

	// Define some vars
	cluster := k.stim.GetConfig("kube-configure-cluster")
	sa := k.stim.GetConfig("kube-service-account")

	// Get secrets from Vault
	vault := k.stim.Vault()
	secretValues, err := vault.GetSecretKeys("secret/kubernetes/" + cluster + "/" + sa + "/kube-config")
	if err != nil {
		return err
	}

	// Set default context
	var context string
	if context = k.stim.GetConfig("kube-context"); context == "" {
		context = cluster
	}

	// Set default namespace
	var namespace string
	if namespace = k.stim.GetConfig("kube-config-namespace"); namespace == "" {
		namespace = secretValues["default-namespace"]
	}

	// Build the config options
	kubeConfigOptions := &kubepkg.KubeConfigOptions{
		ClusterName:             cluster,
		ClusterServer:           secretValues["cluster-server"],
		ClusterCA:               secretValues["cluster-ca"],
		AuthName:                cluster + "-" + sa,
		AuthToken:               secretValues["user-token"],
		ContextName:             context,
		ContextSetCurrent:       k.stim.GetConfigBool("kube-current-context"),
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

// Validates fields for setting up kubeconfig
func (k *Kubernetes) validateConfigureFields() error {

	if k.stim.GetConfig("kube-configure-cluster") == "" {
		return errors.New("Must specify cluster name")
	}

	if k.stim.GetConfig("kube-service-account") == "" {
		return errors.New("Must specify service account")
	}

	return nil
}

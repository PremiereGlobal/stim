package kubernetes

import (
	"github.com/PremiereGlobal/stim/pkg/kubernetes"
	// "github.com/davecgh/go-spew/spew"
)

func (k *Kubernetes) configureContext() error {

	// Create a Vault instance
	k.vault = k.stim.Vault()

	var err error

	cluster, err := k.stim.PromptListVault("secret/kubernetes", "Select Cluster", k.stim.ConfigGetString("kube-config-cluster"), k.stim.ConfigGetString("kube.cluster.filter"))
	if err != nil {
		return err
	}

	sa, err := k.stim.PromptListVault("secret/kubernetes/"+cluster, "Select Service Account", k.stim.ConfigGetString("kube-service-account"), k.stim.ConfigGetString("kube.service-account.filter"))
	if err != nil {
		return err
	}

	// Get secrets from Vault
	secretValues, err := k.vault.GetSecretKeys("secret/kubernetes/" + cluster + "/" + sa + "/kube-config")
	if err != nil {
		return err
	}

	namespace := k.stim.ConfigGetString("kube-config-namespace")
	if namespace == "" {
		namespace, err = k.stim.PromptString("Select Default Namespace", secretValues["default-namespace"])
		if err != nil {
			return err
		}
	}

	context := k.stim.ConfigGetString("kube-context")
	if context == "" {
		context, err = k.stim.PromptString("Context Name", cluster)
		if err != nil {
			return err
		}
	}

	currentContext, err := k.stim.PromptBool("Set as current context?", k.stim.ConfigGetBool("kube-current-context"), true)
	if err != nil {
		return err
	}

	// Build the config options
	kubeConfigOptions := &kubernetes.ConfigOptions{
		ClusterName:             cluster,
		ClusterServer:           secretValues["cluster-server"],
		ClusterCA:               secretValues["cluster-ca"],
		AuthName:                cluster + "-" + sa,
		AuthToken:               secretValues["user-token"],
		ContextName:             context,
		ContextSetCurrent:       currentContext,
		ContextDefaultNamespace: namespace,
	}

	// Gets us a kubeConfig object using the default kubeconfig paths, etc.
	kubeConfig := kubernetes.NewConfig()
	err = kubeConfig.Modify(kubeConfigOptions)
	if err != nil {
		return err
	}

	return nil
}

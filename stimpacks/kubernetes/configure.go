package kubernetes

import (
	"github.com/PremiereGlobal/stim/pkg/kubernetes"
	// "github.com/davecgh/go-spew/spew"
)

func (k *Kubernetes) configureContext() error {

	// Create a Vault instance
	k.vault = k.stim.Vault()

	var err error

	cluster := k.stim.ConfigGetString("kube.config.cluster")
	if cluster == "" {
		cluster = k.stim.ConfigGetString("kube-config-cluster") //TODO: depreciated config should be removed
	}
	kubeClusterFilter := k.stim.ConfigGetString("kube.cluster.filter")

	kubePath := k.stim.ConfigGetString("kube.config.path")

	cluster, err = k.stim.PromptListVault(kubePath, "Select Cluster", cluster, kubeClusterFilter)
	if err != nil {
		return err
	}

	sa := k.stim.ConfigGetString("kube.config.serviceaccount")
	if sa == "" {
		sa = k.stim.ConfigGetString("kube-service-account") //TODO: depreciated config should be removed
	}
	saFilter := k.stim.ConfigGetString("kube.config.serviceaccountfilter")

	sa, err = k.stim.PromptListVault(kubePath+"/"+cluster, "Select Service Account", sa, saFilter)
	if err != nil {
		return err
	}

	kubeKeyName := k.stim.ConfigGetString("kube.config.keyname")

	// Get secrets from Vault
	secretValues, err := k.vault.GetSecretKeys(kubePath + "/" + cluster + "/" + sa + "/" + kubeKeyName)
	if err != nil {
		return err
	}

	namespace := k.stim.ConfigGetString("kube.config.namespace")

	if namespace == "" {
		namespace, err = k.stim.PromptString("Select Default Namespace", secretValues["default-namespace"])
		if err != nil {
			return err
		}
	}

	context := k.stim.ConfigGetString("kube.config.context")
	if context == "" {
		context = k.stim.ConfigGetString("kube-context") //TODO: depreciated config should be removed
	}

	if context == "" {
		context, err = k.stim.PromptString("Context Name", cluster)
		if err != nil {
			return err
		}
	}

	kcc := false
	if k.stim.ConfigHasValue("kube.config.setcontext") {
		kcc = k.stim.ConfigGetBool("kube.config.setcontext")
	} else {
		kcc = k.stim.ConfigGetBool("kube-current-context") //TODO: depreciated config should be removed
	}

	currentContext, err := k.stim.PromptBool("Set as current context?", kcc, true)
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

package kubernetes

import (
	"github.com/PremiereGlobal/stim/pkg/kubernetes"
	"github.com/PremiereGlobal/stim/pkg/utils"

	// "github.com/davecgh/go-spew/spew"
	"sort"
	"strings"
)

func (k *Kubernetes) configureContext() error {

	// Create a Vault instance
	k.vault = k.stim.Vault()

	var err error

	cluster := k.stim.ConfigGetString("kube.config.cluster")
	if cluster == "" {
		cluster = k.stim.ConfigGetString("kube-config-cluster") //TODO: depreciated config should be removed
	}
	kubeClusterFilter := k.stim.ConfigGetString("kube.config.cluster-filter")
	if kubeClusterFilter == "" {
		kubeClusterFilter = k.stim.ConfigGetString("kube.cluster.filter") //TODO: depreciated config should be removed
	}

	kubePath := k.stim.ConfigGetString("kube.config.path")

	cluster, err = k.stim.PromptListVault(kubePath, "Select Cluster", cluster, kubeClusterFilter)
	if err != nil {
		return err
	}

	sa := k.stim.ConfigGetString("kube.config.serviceaccount")
	if sa == "" {
		sa = k.stim.ConfigGetString("kube-service-account") //TODO: depreciated config should be removed
	}
	saFilter := k.stim.ConfigGetString("kube.config.service-account-filter")
	if saFilter == "" {
		saFilter = k.stim.ConfigGetString("kube.config.serviceaccountfilter") //TODO: depreciated config should be removed
	}

	filteredServiceAccounts, err := k.filterServiceAccounts(kubePath+"/"+cluster, saFilter)
	if err != nil {
		return err
	}
	sa, err = k.stim.PromptList("Select Service Account", filteredServiceAccounts, sa)
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

func (k *Kubernetes) filterServiceAccounts(path string, saFilter string) ([]string, error) {
	serviceAccounts, err := k.vault.ListSecrets(path)
	k.stim.Fatal(err)
	regexFilteredServiceAccounts, err := utils.Filter(serviceAccounts, saFilter)
	k.stim.Fatal(err)

	if !k.stim.ConfigGetBool("kube.config.filter-by-token") {
		return regexFilteredServiceAccounts, nil
	}

	kubeKeyName := k.stim.ConfigGetString("kube.config.keyname")
	var paths []string
	for _, serviceAccount := range regexFilteredServiceAccounts {
		paths = append(paths, path+"/"+serviceAccount+"/"+kubeKeyName)
	}

	filteredPaths, err := k.vault.Filter(paths, []string{"read"})
	k.stim.Fatal(err)

	var tokenFilteredServiceAccounts []string
	for _, filteredPath := range filteredPaths {
		sa := strings.TrimSuffix(filteredPath, "/"+kubeKeyName)
		sa = strings.TrimPrefix(sa, path+"/")
		tokenFilteredServiceAccounts = append(tokenFilteredServiceAccounts, sa)
	}

	sort.Strings(tokenFilteredServiceAccounts)
	return tokenFilteredServiceAccounts, nil
}

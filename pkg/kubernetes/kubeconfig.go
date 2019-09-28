package kubernetes

import (
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type KubeConfigOptions struct {
	ClusterName   string
	ClusterServer string
	ClusterCA     string
	AuthName      string
	AuthToken     string
	ContextName   string

	// Set to true to make context current
	ContextSetCurrent bool

	// Default namespace
	ContextDefaultNamespace string

	// Path to an explicit kubeconfig file
	KubeConfigFilePath string
}

type AuthOptions struct {
	name  string
	token string
}

func (k *Kubernetes) modifyKubeconfig(o *KubeConfigOptions) error {

	// GetStartingConfig returns the config that subcommands should being operating against.  It may or may not be merged depending on loading rules
	kubeConfig, err := k.configAccess.GetStartingConfig()
	if err != nil {
		return err
	}

	// Set the kubeconfig components
	k.modifyKubeconfigCluster(kubeConfig, o)
	k.modifyKubeconfigAuth(kubeConfig, o)
	k.modifyKubeconfigContext(kubeConfig, o)

	// Write the config file
	if err := clientcmd.ModifyConfig(k.configAccess, *kubeConfig, true); err != nil {
		return err
	}

	k.log.Debug("Kubernetes config modified...")

	return nil

}

func (k *Kubernetes) modifyKubeconfigCluster(kubeConfig *clientcmdapi.Config, o *KubeConfigOptions) {

	stanza, exists := kubeConfig.Clusters[o.ClusterName]
	if !exists {
		k.log.Debug("Kubernetes cluster `" + o.ClusterName + "` not found in config, creating...")
		stanza = clientcmdapi.NewCluster()
	} else {
		k.log.Debug("Kubernetes cluster `" + o.ClusterName + "` found, modifying...")
	}

	stanza.Server = o.ClusterServer
	stanza.CertificateAuthorityData = []byte(o.ClusterCA)
	stanza.InsecureSkipTLSVerify = false
	stanza.CertificateAuthority = ""
	kubeConfig.Clusters[o.ClusterName] = stanza

}

func (k *Kubernetes) modifyKubeconfigAuth(kubeConfig *clientcmdapi.Config, o *KubeConfigOptions) {

	stanza, exists := kubeConfig.AuthInfos[o.AuthName]
	if !exists {
		k.log.Debug("Kubernetes auth `" + o.AuthName + "` not found in config, creating...")
		stanza = clientcmdapi.NewAuthInfo()
	} else {
		k.log.Debug("Kubernetes auth `" + o.AuthName + "` found, modifying...")
	}

	stanza.Token = o.AuthToken
	kubeConfig.AuthInfos[o.AuthName] = stanza

}

func (k *Kubernetes) modifyKubeconfigContext(kubeConfig *clientcmdapi.Config, o *KubeConfigOptions) {

	stanza, exists := kubeConfig.Contexts[o.ContextName]
	if !exists {
		k.log.Debug("Kubernetes context `" + o.ContextName + "` not found in config, creating...")
		stanza = clientcmdapi.NewContext()
	} else {
		k.log.Debug("Kubernetes context `" + o.ContextName + "` found, modifying...")
	}

	stanza.Cluster = o.ClusterName
	stanza.AuthInfo = o.AuthName
	stanza.Namespace = o.ContextDefaultNamespace
	kubeConfig.Contexts[o.ContextName] = stanza

	if o.ContextSetCurrent {
		kubeConfig.CurrentContext = o.ContextName
		k.log.Debug("Kubernetes current-context set to `" + o.ContextName + "`")
	}

}

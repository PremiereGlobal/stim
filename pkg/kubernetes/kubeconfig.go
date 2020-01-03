package kubernetes

import (
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// KubeConfig represents a single cluster config
type KubeConfig struct {

	// This allows us to read/write the kube config
	// It takes into account KUBECONFIG env var for setting the location
	configAccess clientcmd.ConfigAccess

	// This is the path of the kubeconfig file
	// path string
	// restClientConfig contains the rest configuration for the Kubernetes client
	// For example, Host, Apipath, Auth, etc.
	// restClientConfig *rest.Config
	//
	// kubeConfigPath string
	//
	// // clientset is the set of Kubernetes api clients that can be invoked
	// // For example, core, networking, storage, etc.
	// // This value contains the current working clientset
	// clientSet *kubernetes.Clientset
}

type KubeConfigOptions struct {
	ContextName   string
	ClusterName   string
	ClusterServer string
	ClusterCA     string
	AuthName      string
	AuthToken     string

	// Set to true to make context current
	ContextSetCurrent bool

	// Default namespace
	ContextDefaultNamespace string

	// Path to an explicit kubeconfig file
	// KubeConfigFilePath string
}

func NewKubeConfig(path string) *KubeConfig {

	kc := KubeConfig{}

	if path != "" {
		pathOptions := &clientcmd.PathOptions{}
		pathOptions.LoadingRules = clientcmd.NewDefaultClientConfigLoadingRules()
		pathOptions.LoadingRules.ExplicitPath = path
		kc.configAccess = pathOptions
	} else {
		kc.configAccess = clientcmd.NewDefaultPathOptions()
	}
	// k.log.Debug("Using kubeconfig file: " + k.configAccess.GetDefaultFilename())

	return &kc
}

// func (k *KubeConfig) SetContext() {
//
// }

// func (kc *KubeConfig) SetKubeconfigPath(path string) {
//   kc.path = path
// }
//
// func (k *Kubernetes) SetKubeconfig(kubeConfigOptions *KubeConfigOptions) error {
//
// 	// configAccess is used by subcommands and methods in this package to load and modify the appropriate config files
//
// 	// If we specified an explicit path, use that
//
//
// 	err := k.modifyKubeconfig(kubeConfigOptions)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

// func (kc *KubeConfig) modifyKubeconfigContext(kubeConfig *clientcmdapi.Config, o *KubeConfigOptions) {
//
// 	stanza, exists := kubeConfig.Contexts[o.ContextName]
// 	if !exists {
// 		// kc.log.Debug("Kubernetes context `" + o.ContextName + "` not found in config, creating...")
// 		stanza = clientcmdapi.NewContext()
// 	} else {
// 		// kc.log.Debug("Kubernetes context `" + o.ContextName + "` found, modifying...")
// 	}
//
// 	stanza.Cluster = o.ClusterName
// 	stanza.AuthInfo = o.AuthName
// 	stanza.Namespace = o.ContextDefaultNamespace
// 	kubeConfig.Contexts[o.ContextName] = stanza
//
// 	if o.ContextSetCurrent {
// 		kubeConfig.CurrentContext = o.ContextName
// 		// kc.log.Debug("Kubernetes current-context set to `" + o.ContextName + "`")
// 	}
//
// }

func (kc *KubeConfig) ModifyConfig(o *KubeConfigOptions) error {

	// GetStartingConfig returns the config that subcommands should being operating against.  It may or may not be merged depending on loading rules
	kubeConfig, err := kc.configAccess.GetStartingConfig()
	if err != nil {
		return err
	}

	// Set the kubeconfig components
	kc.modifyKubeconfigCluster(kubeConfig, o)
	kc.modifyKubeconfigAuth(kubeConfig, o)
	kc.modifyKubeconfigContext(kubeConfig, o)

	// Write the config file
	if err := clientcmd.ModifyConfig(kc.configAccess, *kubeConfig, true); err != nil {
		return err
	}

	// kc.log.Debug("Kubernetes config modified...")

	return nil

}

func (kc *KubeConfig) modifyKubeconfigCluster(kubeConfig *clientcmdapi.Config, o *KubeConfigOptions) {

	stanza, exists := kubeConfig.Clusters[o.ClusterName]
	if !exists {
		// kc.log.Debug("Kubernetes cluster `" + o.ClusterName + "` not found in config, creating...")
		stanza = clientcmdapi.NewCluster()
	} else {
		// kc.log.Debug("Kubernetes cluster `" + o.ClusterName + "` found, modifying...")
	}

	stanza.Server = o.ClusterServer
	stanza.CertificateAuthorityData = []byte(o.ClusterCA)
	stanza.InsecureSkipTLSVerify = false
	stanza.CertificateAuthority = ""
	kubeConfig.Clusters[o.ClusterName] = stanza

}

func (kc *KubeConfig) modifyKubeconfigAuth(kubeConfig *clientcmdapi.Config, o *KubeConfigOptions) {

	stanza, exists := kubeConfig.AuthInfos[o.AuthName]
	if !exists {
		// kc.log.Debug("Kubernetes auth `" + o.AuthName + "` not found in config, creating...")
		stanza = clientcmdapi.NewAuthInfo()
	} else {
		// kc.log.Debug("Kubernetes auth `" + o.AuthName + "` found, modifying...")
	}

	stanza.Token = o.AuthToken
	kubeConfig.AuthInfos[o.AuthName] = stanza

}

func (kc *KubeConfig) modifyKubeconfigContext(kubeConfig *clientcmdapi.Config, o *KubeConfigOptions) {

	stanza, exists := kubeConfig.Contexts[o.ContextName]
	if !exists {
		// kc.log.Debug("Kubernetes context `" + o.ContextName + "` not found in config, creating...")
		stanza = clientcmdapi.NewContext()
	}
	// } else {
	// 	// kc.log.Debug("Kubernetes context `" + o.ContextName + "` found, modifying...")
	// }

	stanza.Cluster = o.ClusterName
	stanza.AuthInfo = o.AuthName
	stanza.Namespace = o.ContextDefaultNamespace
	kubeConfig.Contexts[o.ContextName] = stanza

	if o.ContextSetCurrent {
		kubeConfig.CurrentContext = o.ContextName
		// kc.log.Debug("Kubernetes current-context set to `" + o.ContextName + "`")
	}

}

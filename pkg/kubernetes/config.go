package kubernetes

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// Config effectivly represents a kubeconfig file configuration allowing for
// creating or modifying the values in that file
type Config struct {

	// configAccess contains the configuration for which kubeconfig file(s) we're
	// dealing with
	configAccess clientcmd.ConfigAccess
}

// ConfigOptions defines options for configuring the kubeconfig
type ConfigOptions struct {

	// ClusterName is the name of the cluster as it appears in the kubeconfig
	ClusterName string

	// ClusterServer is the hostname of the cluster
	ClusterServer string

	// ClusterCA is the Certificate Authority contents
	ClusterCA string

	// AuthName is the name of the authentication entry as it appears in the kubeconfig
	AuthName string

	// AuthToken is the authentication token
	AuthToken string

	// ContextName is the name of the context
	ContextName string

	// ContextDefaultNamespace is the context's default namespace
	ContextDefaultNamespace string

	// ContextSetCurrent is a flag to set this context as the "current-context" in the kubeconfig
	ContextSetCurrent bool
}

// NewConfig creates a new config object using the environment defaults
func NewConfig() *Config {

	config := &Config{}
	config.configAccess = clientcmd.NewDefaultPathOptions()

	return config
}

// NewConfigFromPath creates a new config object using the specified kubeconfig path
func NewConfigFromPath(kubeConfigFilePath string) *Config {

	config := &Config{}

	pathOptions := &clientcmd.PathOptions{}
	pathOptions.LoadingRules = clientcmd.NewDefaultClientConfigLoadingRules()
	pathOptions.LoadingRules.ExplicitPath = kubeConfigFilePath
	config.configAccess = pathOptions

	return config
}

// Modify updates the kubeconfig with the given options
func (c *Config) Modify(options *ConfigOptions) error {

	newConfig, err := c.configAccess.GetStartingConfig()
	if err != nil {
		return err
	}

	cluster := clientcmdapi.NewCluster()
	cluster.Server = options.ClusterServer
	cluster.CertificateAuthorityData = []byte(options.ClusterCA)
	newConfig.Clusters[options.ClusterName] = cluster

	authInfo := clientcmdapi.NewAuthInfo()
	authInfo.Token = options.AuthToken
	newConfig.AuthInfos[options.AuthName] = authInfo

	context := clientcmdapi.NewContext()
	if options.ContextDefaultNamespace != "" {
		context.Namespace = options.ContextDefaultNamespace
	}
	context.Cluster = options.ClusterName
	context.AuthInfo = options.AuthName
	newConfig.Contexts[options.ContextName] = context

	if options.ContextSetCurrent {
		newConfig.CurrentContext = options.ContextName
	}

	clientcmd.ModifyConfig(c.configAccess, *newConfig, false)

	return nil
}

// GetRestClientConfig returns a rest.Config to be used in a Kubernetes client
func (c *Config) GetRestClientConfig() (*rest.Config, error) {

	// This loads in the kubeconfig file
	clientcmdapiConfig, err := c.configAccess.GetStartingConfig()
	if err != nil {
		return nil, err
	}

	// Creates the Kubernetes client based on the given kubeconfig
	clientConfig, err := clientcmd.NewDefaultClientConfig(*clientcmdapiConfig, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, err
	}

	return clientConfig, nil
}

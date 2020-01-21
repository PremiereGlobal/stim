package kubernetes

import (
	"k8s.io/client-go/kubernetes"
)

// Kubernetes represents a interface with a Kubernetes cluster
type Kubernetes struct {
	config *Config
}

// New returns a new Kubernetes object with the given config
func New(config *Config) (*Kubernetes, error) {
	return &Kubernetes{config: config}, nil
}

// GetConfig returns the config object assicated with this instance
func (k *Kubernetes) GetConfig() *Config {
	return k.config
}

// GetClientset returns the Kubernetes clientset, which can be acted on directly
// to perform API calls
func (k *Kubernetes) GetClientset() (*kubernetes.Clientset, error) {

	restClientConfig, err := k.GetConfig().GetRestClientConfig()
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(restClientConfig)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}

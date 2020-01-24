package kubernetes

import (
	"k8s.io/client-go/discovery"
)

// DiscoveryClient returns the Kubernetes discovery client
// The discovery client implements the functions that discover server-supported API groups,
// versions and resources.
func (k *Kubernetes) DiscoveryClient() (*discovery.DiscoveryClient, error) {

	restClientConfig, err := k.GetConfig().GetRestClientConfig()
	if err != nil {
		return nil, err
	}

	discoveryClient, err := discovery.NewDiscoveryClientForConfig(restClientConfig)
	if err != nil {
		return nil, err
	}

	return discoveryClient, nil
}

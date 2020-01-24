package kubernetes

// Version returns the version of the Kubernetes cluster
func (k *Kubernetes) Version() (string, error) {

	discoveryClient, err := k.DiscoveryClient()
	if err != nil {
		return "", err
	}

	serverVersion, err := discoveryClient.ServerVersion()
	if err != nil {
		return "", err
	}

	return serverVersion.GitVersion, nil
}

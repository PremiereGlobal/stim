package kubernetes

import (
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
)

type Kubernetes struct {
	// client *api.Client
	config *Config

	// This allows us to read/write the kube config
	// It takes into account KUBECONFIG env var for setting the location
	configAccess clientcmd.ConfigAccess
}

type Config struct {
	// Address string
	Logger
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
}

func (k *Kubernetes) Debug(message string) {
	if k.config.Logger != nil {
		k.config.Debug(message)
	}
}

func (k *Kubernetes) Info(message string) {
	if k.config.Logger != nil {
		k.config.Info(message)
	} else {
		fmt.Println(message)
	}
}

func New(config *Config) (*Kubernetes, error) {

	k := &Kubernetes{config: config}

	return k, nil
}

func (k *Kubernetes) SetKubeconfig(kubeConfigOptions *KubeConfigOptions) error {

	// configAccess is used by subcommands and methods in this package to load and modify the appropriate config files
	k.configAccess = clientcmd.NewDefaultPathOptions()
	k.Debug("Using kubeconfig file: " + k.configAccess.GetDefaultFilename())

	err := k.modifyKubeconfig(kubeConfigOptions)
	if err != nil {
		return err
	}

	return nil
}

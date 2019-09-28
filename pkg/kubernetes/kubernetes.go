package kubernetes

import (
	"github.com/PremiereGlobal/stim/pkg/stimlog"
	"k8s.io/client-go/tools/clientcmd"
)

type Kubernetes struct {
	// client *api.Client
	config *Config
	log    Logger
	// This allows us to read/write the kube config
	// It takes into account KUBECONFIG env var for setting the location
	configAccess clientcmd.ConfigAccess
}

type Config struct {
	Log Logger
}

type Logger interface {
	Debug(...interface{})
	Warn(...interface{})
	Fatal(...interface{})
}

func New(kconf *Config) (*Kubernetes, error) {

	k := &Kubernetes{config: kconf}

	if kconf.Log != nil {
		k.log = kconf.Log
	} else {
		k.log = stimlog.GetLogger()
	}

	return k, nil
}

func (k *Kubernetes) SetKubeconfig(kubeConfigOptions *KubeConfigOptions) error {

	// configAccess is used by subcommands and methods in this package to load and modify the appropriate config files
	// If we specified an explicit path, use that
	if kubeConfigOptions.KubeConfigFilePath != "" {
		pathOptions := &clientcmd.PathOptions{}
		pathOptions.LoadingRules = clientcmd.NewDefaultClientConfigLoadingRules()
		pathOptions.LoadingRules.ExplicitPath = kubeConfigOptions.KubeConfigFilePath
		k.configAccess = pathOptions
	} else {
		k.configAccess = clientcmd.NewDefaultPathOptions()
	}
	k.log.Debug("Using kubeconfig file: " + k.configAccess.GetDefaultFilename())

	err := k.modifyKubeconfig(kubeConfigOptions)
	if err != nil {
		return err
	}

	return nil
}

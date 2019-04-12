package kubernetes

import (
	"github.com/readytalk/stim/pkg/stimlog"
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
	// Address string
}

type Logger interface {
	Debug(...interface{})
	Warn(...interface{})
	Fatal(...interface{})
}

func New(kconf *Config, log Logger) (*Kubernetes, error) {

	k := &Kubernetes{config: kconf}

	if log != nil {
		k.log = log
	} else {
		k.log = stimlog.GetLogger()
	}

	return k, nil
}

func (k *Kubernetes) SetKubeconfig(kubeConfigOptions *KubeConfigOptions) error {

	// configAccess is used by subcommands and methods in this package to load and modify the appropriate config files
	k.configAccess = clientcmd.NewDefaultPathOptions()
	k.log.Debug("Using kubeconfig file: " + k.configAccess.GetDefaultFilename())

	err := k.modifyKubeconfig(kubeConfigOptions)
	if err != nil {
		return err
	}

	return nil
}

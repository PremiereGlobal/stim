package kubernetes

import (
	// "fmt"
	// "errors"
	"github.com/PremiereGlobal/stim/pkg/stimlog"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Kubernetes represents a connection to a Kubernetes cluster
type Kubernetes struct {
	// client *api.Client
	config *Config
	log    Logger

	// This allows us to read/write the kube config
	// It takes into account KUBECONFIG env var for setting the location
	configAccess clientcmd.ConfigAccess

	// restClientConfig contains the rest configuration for the Kubernetes client
	// For example, Host, Apipath, Auth, etc.
	// restClientConfig *rest.Config

	kubeConfigPath string

	// clientset is the set of Kubernetes api clients that can be invoked
	// For example, core, networking, storage, etc.
	// This value contains the current working clientset
	clientSet *kubernetes.Clientset
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

	// Make a new config based on our env
	// kubeConfig := clientcmd.NewDefaultPathOptions()

	// fmt.Println(kubeConfig.GetDefaultFilename())
	// return nil, nil
	// kubeConfigPath := filepath.Join(home, ".kube", "config")
	// if e, _ := exists(kubeConfigPath); home != "" && e {
	//   log.Println("Attempting .kube/config")
	//   // use the current context in kubeconfig
	//   kubeConfig, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	//   if err != nil {
	//     log.Fatal(err.Error())
	//   }
	// }

	// Make a new Kubernetes client based on our current config
	// k8sclient, err := kubernetes.NewForConfig(kubeConfig)
	// if err != nil {
	//   return nil, err
	// }

	k := &Kubernetes{config: kconf}

	if kconf.Log != nil {
		k.log = kconf.Log
	} else {
		k.log = stimlog.GetLogger()
	}

	return k, nil
}

// SetKubeconfig updates the working config based on a provided kubeconfig file path
func (k *Kubernetes) SetKubeconfigPath(kubeconfigPath string) error {

	// Create the new restClientConfig based upon the provide kubeconfig file
	restClientConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return err
	}

	// Update our clientSet based upon the new configuration
	err = k.updateClientset(restClientConfig)
	if err != nil {
		return err
	}

	return nil
}

//
// func GetKubeconfigPath() {
//
// 	// Create the new restClientConfig based upon the provide kubeconfig file
// 	restClientConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
// 	if err != nil {
// 		return error
// 	}
//
// 	// Update our clientSet based upon the new configuration
// 	err = k.updateClientset(restClientConfig)
// 	if err != nil {
// 		return error
// 	}
//
// 	return nil
// }

// updateClientset updates the working clientset with the provided restClientConfig
func (k *Kubernetes) updateClientset(restClientConfig *rest.Config) error {

	clientSet, err := kubernetes.NewForConfig(restClientConfig)
	if err != nil {
		return err
	}
	k.clientSet = clientSet

	return nil
}

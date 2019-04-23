package stim

import (
	"github.com/PremiereGlobal/stim/pkg/kubernetes"
)

func (stim *Stim) Kubernetes() *kubernetes.Kubernetes {
	stim.log.Debug("Stim-Kubernetes: Creating")

	k, err := kubernetes.New(&kubernetes.Config{Log: stim.log})
	if err != nil {
		stim.log.Fatal("Stim-Kubernetes: Error Initializaing: ", err)
	}

	return k
}

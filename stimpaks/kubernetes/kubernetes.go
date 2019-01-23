package kubernetes

import (
	"github.com/readytalk/stim/stim"
)

type Kubernetes struct {
	name string
	stim *stim.Stim
}

func New() *Kubernetes {
	k := &Kubernetes{name: "kubernetes"}
	return k
}

func (k *Kubernetes) Name() string {
	return k.name
}

func (k *Kubernetes) BindStim(s *stim.Stim) {
	k.stim = s
}

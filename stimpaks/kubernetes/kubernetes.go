package kubernetes

import (
	vault "github.com/readytalk/stim/pkg/vault"
	"github.com/readytalk/stim/stim"
)

type Kubernetes struct {
	name  string
	stim  *stim.Stim
	vault *vault.Vault
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

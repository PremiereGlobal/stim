package vault

import (
	"github.com/PremiereGlobal/stim/stim"
)

type Vault struct {
	name string
	stim *stim.Stim
}

func New() *Vault {
	vault := &Vault{name: "vault"}
	return vault
}

func (v *Vault) Name() string {
	return v.name
}

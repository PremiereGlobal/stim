package vault

import (
	"github.com/readytalk/stim/stim"
)

type Vault struct {
	stim *stim.Stim
}

func New() *Vault {
	vault := &Vault{}
	return vault
}

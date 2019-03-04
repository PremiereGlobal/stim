package aws

import (
	"github.com/readytalk/stim/pkg/log"
	vault "github.com/readytalk/stim/pkg/vault"
	"github.com/readytalk/stim/stim"
)

type Aws struct {
	name  string
	stim  *stim.Stim
	vault *vault.Vault
	Log   log.Logger
}

func New() *Aws {
	aws := &Aws{name: "aws"}
	log.SetLogger(aws.Log)
	return aws
}

func (a *Aws) Name() string {
	return a.name
}

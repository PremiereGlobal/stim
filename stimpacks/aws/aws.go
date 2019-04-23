package aws

import (
	"github.com/PremiereGlobal/stim/pkg/stimlog"
	vault "github.com/PremiereGlobal/stim/pkg/vault"
	"github.com/PremiereGlobal/stim/stim"
)

type Aws struct {
	name  string
	stim  *stim.Stim
	vault *vault.Vault
	log   stimlog.StimLogger
}

func New() *Aws {
	aws := &Aws{name: "aws"}
	aws.log = stimlog.GetLogger()
	return aws
}

func (a *Aws) Name() string {
	return a.name
}

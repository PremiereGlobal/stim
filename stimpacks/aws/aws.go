package aws

import (
	"github.com/readytalk/stim/pkg/stimlog"
	vault "github.com/readytalk/stim/pkg/vault"
	"github.com/readytalk/stim/stim"
)

type Aws struct {
	name  string
	stim  *stim.Stim
	vault *vault.Vault
	log   *stimlog.StimLogger
}

func New() *Aws {
	aws := &Aws{name: "aws"}
	aws.log = stimlog.GetLogger()
	return aws
}

func (a *Aws) Name() string {
	return a.name
}

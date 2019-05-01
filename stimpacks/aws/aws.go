package aws

import (
	stimaws "github.com/PremiereGlobal/stim/pkg/aws"
	"github.com/PremiereGlobal/stim/pkg/stimlog"
	"github.com/PremiereGlobal/stim/pkg/vault"
	"github.com/PremiereGlobal/stim/stim"
)

type Aws struct {
	name  string
	stim  *stim.Stim
	vault *vault.Vault
	aws   *stimaws.Aws
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

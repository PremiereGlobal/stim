package aws

import (
	"github.com/readytalk/stim/stim"
)

type Aws struct {
	name string
	stim *stim.Stim
}

func New() *Aws {
	aws := &Aws{name: "aws"}
	return aws
}

func (a *Aws) Name() string {
	return a.name
}

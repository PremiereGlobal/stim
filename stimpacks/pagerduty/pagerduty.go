package pagerduty

import (
	"github.com/PremiereGlobal/stim/stim"
)

type Pagerduty struct {
	name string
	stim *stim.Stim
}

func New() *Pagerduty {
	pagerduty := &Pagerduty{name: "pagerduty"}
	return pagerduty
}

func (p *Pagerduty) Name() string {
	return p.name
}

package pagerduty

import (
	pd "github.com/readytalk/stim/pkg/pagerduty"
	"github.com/readytalk/stim/stim"
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

func (p *Pagerduty) SendEvent() {

	pagerduty := p.stim.Pagerduty()

	err := pagerduty.SendEvent(&pd.Event{
		Service:  p.stim.GetConfig("pagerduty-service"),
		Summary:  p.stim.GetConfig("pagerduty-summary"),
		Action:   p.stim.GetConfig("pagerduty-action"),
		Severity: p.stim.GetConfig("pagerduty-severity"),
	})
	p.stim.Fatal(err)

}

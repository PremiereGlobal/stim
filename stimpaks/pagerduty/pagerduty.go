package pagerduty

import (
	"github.com/readytalk/stim/api"
	pd "github.com/readytalk/stim/pkg/pagerduty"
	"github.com/spf13/cobra"
)

type Pagerduty struct {
	api *api.Api
	log api.Log
}

func New(api *api.Api) *Pagerduty {
	pagerduty := &Pagerduty{api: api, log: api.Log}
	return pagerduty
}

func (p *Pagerduty) Bind(parentCmd *cobra.Command) {
	p.BindCommand(parentCmd)
}

func (p *Pagerduty) SendEvent() {

	pagerduty := p.api.Pagerduty()

	p.log.Info("Sending event to Pagerduty")
	err := pagerduty.SendEvent(&pd.Event{Summary: "b"})
	if err != nil {
		p.log.Fatal(err)
	}
}

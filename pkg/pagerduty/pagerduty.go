package pagerduty

import (
  "os"
	pdApi "github.com/PagerDuty/go-pagerduty"
	"github.com/readytalk/stim/pkg/utils"
	VaultApi "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
)

type Pagerduty struct {
	client  *pdApi.Client
	log     *logrus.Logger
	vault   *VaultApi.Client
	Options *Options
}

type Options struct {
	Action          string `mapstructure:"notify-pagerduty-action"`
	Service         string `mapstructure:"notify-pagerduty-service"`
	Severity        string `mapstructure:"notify-pagerduty-severity"`
	Summary         string `mapstructure:"notify-pagerduty-summary"`
	Source          string `mapstructure:"notify-pagerduty-source"`
	Component       string `mapstructure:"notify-pagerduty-component"`
	Group           string `mapstructure:"notify-pagerduty-group"`
	Class           string `mapstructure:"notify-pagerduty-class"`
	Details         string `mapstructure:"notify-pagerduty-details"`
	DedupKey        string `mapstructure:"notify-pagerduty-dedupkey"`
	VaultApikeyPath string `mapstructure:"vault-pagerduty-apikey-path"`
	VaultApikeyKey  string `mapstructure:"vault-pagerduty-apikey-key"`
}

func New() *Pagerduty {

	o := &Options{}
	p := &Pagerduty{Options: o}
	return p
}

func (p *Pagerduty) Init() {

	p.log.Info("Fetching Pagerduty API key from Vault...")
	secret, err := p.vault.Logical().Read(p.Options.VaultApikeyPath)
	if err != nil {
		p.log.Fatal(err)
	}
	apikey := secret.Data[p.Options.VaultApikeyKey].(string)
	client := pdApi.NewClient(apikey)
	p.client = client
}

func (p *Pagerduty) GetServiceIntegrationId(servicename string) string {

	var serviceid string
	var integrationid string

	opts := pdApi.ListServiceOptions{Query: servicename, Includes: []string{"integrations"}}

	if svcs, err := p.client.ListServices(opts); err != nil {
		panic(err)
	} else {
		for _, s := range svcs.Services {
			if s.Name == servicename {
				serviceid = s.ID
				p.log.Debug("Pagerduty service \"", servicename, "\" found")
				for _, i := range s.Integrations {
					if i.Type == "events_api_v2_inbound_integration" {
						p.log.Debug("Pagerduty EventsAPIv2 integration ID found")
						integrationid = i.IntegrationKey
					}
				}
			}
		}
	}

	if serviceid == "" {
		p.log.Fatal("Pagerduty service \"", servicename, "\" not found")
	}

	if integrationid == "" {
		p.log.Debug("Pagerduty EventsAPIv2 integration ID not found, creating it")
		var integration pdApi.Integration
		integration.Type = "events_api_v2_inbound_integration"
		i, err := p.client.CreateIntegration(serviceid, integration)
		if err != nil {
			p.log.Fatal("Error creating EventsAPIv2 integration ID for service ", err)
		}
		integrationid = i.IntegrationKey
	}

	return integrationid
}

func (p *Pagerduty) SendEvent() {

	o := p.Options

	integrationid := p.GetServiceIntegrationId(o.Service)

  source, _ := os.Hostname()

	payload := pdApi.V2Payload{
		Summary:   o.Summary,
		Source:    source,
		Severity:  o.Severity,
		Component: o.Component,
		Group:     o.Group,
		Class:     o.Class,
		Details:   o.Details,
	}

	event := pdApi.V2Event{
		RoutingKey: integrationid,
		Action:     o.Action,
		Payload:    &payload,
		DedupKey:   o.DedupKey,
	}

	p.log.Info("Sending event to Pagerduty")
	response, err := pdApi.ManageEvent(event)
	if err != nil {
		p.log.Fatal(err)
	}

	p.log.Debug(response)

}

func (p *Pagerduty) CheckRequiredFields() {
	if p.Options.Service == "" {
		p.log.Fatal("Service name (--service) must be set")
	}
	if p.Options.Summary == "" {
		p.log.Fatal("Summary (--summary) must be set")
	}
	if p.Options.Action == "" {
		p.log.Fatal("Action (--action) must be set")
	}
	actions := []string{"trigger", "acknowledge", "resolve"}
	if !utils.Contains(actions, p.Options.Action) {
		p.log.Fatal("Invalid value for action (--action)")
	}
	if p.Options.Severity == "" {
		p.log.Fatal("Severity (--severity) must be set")
	}
	severities := []string{"critical", "warning", "error", "info"}
	if !utils.Contains(severities, p.Options.Severity) {
		p.log.Fatal("Invalid value for severity (--severity)")
	}
	if p.Options.Source == "" {
		p.log.Fatal("Source (--source) must be set")
	}
	if p.Options.VaultApikeyPath == "" {
		p.log.Fatal("Vault APIkey path for Pagerduty not found")
	}
	if p.Options.VaultApikeyKey == "" {
		p.log.Fatal("Vault APIkey key for Pagerduty not found")
	}
	if p.Options.Summary == "" {
		p.log.Fatal("Summary (--summary) must be set")
	}
}

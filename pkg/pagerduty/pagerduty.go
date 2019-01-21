package pagerduty

import (
	"errors"
	pdApi "github.com/PagerDuty/go-pagerduty"
	"github.com/readytalk/stim/pkg/utils"
	"os"
	"strings"
)

type Pagerduty struct {
	client *pdApi.Client
}

type Event struct {
	Action    string
	Service   string
	Severity  string
	Summary   string
	Source    string
	Component string
	Group     string
	Class     string
	Details   string
	DedupKey  string
}

func New(apiKey string) *Pagerduty {

	// Initialize client
	client := pdApi.NewClient(apiKey)
	p := &Pagerduty{client: client}

	return p
}

func (p *Pagerduty) GetServiceIntegrationId(servicename string) (string, error) {

	var serviceid string
	var integrationid string

	opts := pdApi.ListServiceOptions{Query: servicename, Includes: []string{"integrations"}}

	if svcs, err := p.client.ListServices(opts); err != nil {
		panic(err)
	} else {
		for _, s := range svcs.Services {
			if s.Name == servicename {
				serviceid = s.ID
				for _, i := range s.Integrations {
					if i.Type == "events_api_v2_inbound_integration" {
						integrationid = i.IntegrationKey
					}
				}
			}
		}
	}

	if serviceid == "" {
		return "", errors.New("Pagerduty service \"" + servicename + "\" not found")
	}

	// Create integration if it doesn't exist
	if integrationid == "" {
		var integration pdApi.Integration
		integration.Type = "events_api_v2_inbound_integration"
		i, err := p.client.CreateIntegration(serviceid, integration)
		if err != nil {
			return "", err
		}
		integrationid = i.IntegrationKey
	}

	return integrationid, nil
}

func (p *Pagerduty) SendEvent(e *Event) error {

	err := p.validateEventFields(e)
	if err != nil {
		return err
	}

	integrationid, err := p.GetServiceIntegrationId(e.Service)
	if err != nil {
		return err
	}

	var source string
	if e.Source == "" {
		source, err = os.Hostname()
		if err != nil {
			source = "unknown"
		}
	} else {
		source = e.Source
	}

	payload := pdApi.V2Payload{
		Summary:   e.Summary,
		Source:    source,
		Severity:  e.Severity,
		Component: e.Component,
		Group:     e.Group,
		Class:     e.Class,
		Details:   e.Details,
	}

	event := pdApi.V2Event{
		RoutingKey: integrationid,
		Action:     e.Action,
		Payload:    &payload,
		DedupKey:   e.DedupKey,
	}

	_, err = pdApi.ManageEvent(event)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pagerduty) validateEventFields(e *Event) error {
	if e.Service == "" {
		return errors.New("Pagerduty: Event Service Name must be set")
	}
	if e.Summary == "" {
		return errors.New("Pagerduty: Event Summary must be set")
	}
	if e.Action == "" {
		return errors.New("Pagerduty: Event Action must be set")
	}
	validActions := []string{"trigger", "acknowledge", "resolve"}
	if !utils.Contains(validActions, e.Action) {
		return errors.New("Pagerduty: Invalid value for Event Action. Valid values are: [" + strings.Join(validActions, ",") + "]")
	}
	if e.Severity == "" {
		return errors.New("Pagerduty: Event Severity must be set")
	}
	validSeverities := []string{"critical", "warning", "error", "info"}
	if !utils.Contains(validSeverities, e.Severity) {
		return errors.New("Pagerduty: Invalid value for Event Severity. Valid values are: [" + strings.Join(validSeverities, ",") + "]")
	}
	if e.Summary == "" {
		return errors.New("Pagerduty: Event Summary must be set")
	}

	return nil
}

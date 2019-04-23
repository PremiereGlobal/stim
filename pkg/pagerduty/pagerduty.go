package pagerduty

import (
	"errors"
	"os"
	"strings"

	pdApi "github.com/PagerDuty/go-pagerduty"
	"github.com/PremiereGlobal/stim/pkg/utils"
)

// Pagerduty is the main object
type Pagerduty struct {
	client *pdApi.Client
	log    Logger
}

// Event contains the required and optional fields to sent an event
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

type Logger interface {
	Debug(...interface{})
	Warn(...interface{})
	Fatal(...interface{})
}

// New returns a new Pagerduty "instance"
func New(apiKey string, log Logger) *Pagerduty {

	// Initialize client
	client := pdApi.NewClient(apiKey)
	p := &Pagerduty{client: client, log: log}

	return p
}

// GetServices returns a list of all services names within the Pagerduty account
func (p *Pagerduty) GetServices() ([]string, error) {

	// `limit` is how many items will be fetched per request
	limit := uint(50)

	// `done` is set when the last page is fetched
	done := false

	// Options for the request
	options := pdApi.ListServiceOptions{APIListObject: pdApi.APIListObject{Offset: 0, Limit: limit}}

	// Initialize results string
	var results []string

	// Loop through pages and get all services
	for done == false {

		services, err := p.client.ListServices(options)
		if err != nil {
			return nil, err
		}

		for _, s := range services.Services {
			results = append(results, s.Name)
		}

		// Bail if this is the last page
		if !services.APIListObject.More {
			done = true
		}

		// Bump the offset for the next run
		options.APIListObject.Offset = options.APIListObject.Offset + limit
	}

	return results, nil
}

// SendEvent sends the provided Event to Pagerduty.  It automatically detects
// and sets the hostname as the `source`, if not set.
func (p *Pagerduty) SendEvent(e *Event) error {

	err := p.validateEventFields(e)
	if err != nil {
		return err
	}

	integrationid, err := p.getServiceIntegrationID(e.Service)
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

func (p *Pagerduty) getServiceIntegrationID(servicename string) (string, error) {

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

package pagerduty

import (
	"errors"
	pd "github.com/PremiereGlobal/stim/pkg/pagerduty"
)

func (p *Pagerduty) SendEvent() {

	var err error

	pagerduty := p.stim.Pagerduty()

	// Prompt for the service name (if not provided)
	serviceName := p.stim.GetConfig("pagerduty-service")
	if serviceName == "" && p.stim.IsAutomated() {
		p.stim.Fatal(errors.New("Pagerduty `service name` not specified"))
	} else if serviceName == "" {

		// Get the channel list
		services, err := pagerduty.GetServices()
		p.stim.Fatal(err)

		// Prompt for channel
		serviceName, err = p.stim.PromptSearchList("Choose Service:", services)
		p.stim.Fatal(err)

	}

	// Prompt for the summary text (if not provided)
	summary := p.stim.GetConfig("pagerduty-summary")
	if summary == "" && p.stim.IsAutomated() {
		p.stim.Fatal(errors.New("Pagerduty `summary` not specified"))
	} else if summary == "" {

		// Prompt
		summary, err = p.stim.PromptString("Summary Text", "")
		p.stim.Fatal(err)

	}

	// Prompt for the action (if not provided)
	action := p.stim.GetConfig("pagerduty-action")
	if action == "" && p.stim.IsAutomated() {
		p.stim.Fatal(errors.New("Pagerduty `action` not specified"))
	} else if action == "" {

		// Prompt
		action, err = p.stim.PromptString("Action (trigger/resolve)", "trigger")
		p.stim.Fatal(err)

	}

	// Prompt for the severity (if not provided)
	severity := p.stim.GetConfig("pagerduty-severity")
	if severity == "" && p.stim.IsAutomated() {
		p.stim.Fatal(errors.New("Pagerduty `severity` not specified"))
	} else if severity == "" {

		// Prompt
		severity, err = p.stim.PromptString("Severity (info, warning, error, critical)", "warning")
		p.stim.Fatal(err)

	}

	event := &pd.Event{
		Service:  serviceName,
		Summary:  summary,
		Action:   action,
		Severity: severity,
	}

	err = pagerduty.SendEvent(event)
	p.stim.Fatal(err)

}

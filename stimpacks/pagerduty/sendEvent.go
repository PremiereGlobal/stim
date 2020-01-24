package pagerduty

import (
	"errors"

	pd "github.com/PremiereGlobal/stim/pkg/pagerduty"
)

func (p *Pagerduty) SendEvent() {

	var err error

	pagerduty := p.stim.Pagerduty()

	// Prompt for the service name (if not provided)
	serviceName := p.stim.ConfigGetString("pagerduty-service")
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
	summary := p.stim.ConfigGetString("pagerduty-summary")
	if summary == "" && p.stim.IsAutomated() {
		p.stim.Fatal(errors.New("Pagerduty `summary` not specified"))
	} else if summary == "" {

		// Prompt
		summary, err = p.stim.PromptString("Summary Text", "")
		p.stim.Fatal(err)

	}

	// Prompt for the detail text (if not provided)
	details := p.stim.ConfigGetString("pagerduty-details")
	if details == "" && !p.stim.IsAutomated() {
		details, err = p.stim.PromptString("Detailed Text", "")
	}

	// Prompt for the source text (if not provided)
	source := p.stim.ConfigGetString("pagerduty-source")
	if source == "" && !p.stim.IsAutomated() {
		source, err = p.stim.PromptString("Source Text", "")
	}

	// Prompt for the component text (if not provided)
	component := p.stim.ConfigGetString("pagerduty-component")
	if component == "" && !p.stim.IsAutomated() {
		component, err = p.stim.PromptString("Component Text", "")
	}

	// Prompt for the group text (if not provided)
	group := p.stim.ConfigGetString("pagerduty-group")
	if group == "" && !p.stim.IsAutomated() {
		group, err = p.stim.PromptString("Group Text", "")
	}

	// Prompt for the class text (if not provided)
	class := p.stim.ConfigGetString("pagerduty-class")
	if class == "" && !p.stim.IsAutomated() {
		class, err = p.stim.PromptString("Class Text", "")
	}

	// Prompt for the dedupkey text (if not provided)
	dedupKey := p.stim.ConfigGetString("pagerduty-dedupkey")
	if dedupKey == "" && !p.stim.IsAutomated() {
		dedupKey, err = p.stim.PromptString("Dedupkey Text", "")
	}

	// Prompt for the action (if not provided)
	action := p.stim.ConfigGetString("pagerduty-action")
	if action == "" && p.stim.IsAutomated() {
		p.stim.Fatal(errors.New("Pagerduty `action` not specified"))
	} else if action == "" {

		// Prompt
		action, err = p.stim.PromptString("Action (trigger/resolve)", "trigger")
		p.stim.Fatal(err)

	}

	// Prompt for the severity (if not provided)
	severity := p.stim.ConfigGetString("pagerduty-severity")
	if severity == "" && p.stim.IsAutomated() {
		p.stim.Fatal(errors.New("Pagerduty `severity` not specified"))
	} else if severity == "" {

		// Prompt
		severity, err = p.stim.PromptString("Severity (info, warning, error, critical)", "warning")
		p.stim.Fatal(err)

	}

	event := &pd.Event{
		Service:   serviceName,
		Summary:   summary,
		Details:   details,
		Action:    action,
		Severity:  severity,
		Source:    source,
		Component: component,
		Group:     group,
		Class:     class,
		DedupKey:  dedupKey,
	}

	err = pagerduty.SendEvent(event)
	p.stim.Fatal(err)

}

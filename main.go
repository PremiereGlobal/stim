package main

import (
	"github.com/readytalk/stim/stim"
	"github.com/readytalk/stim/stimpacks/aws"
	"github.com/readytalk/stim/stimpacks/discover"
	"github.com/readytalk/stim/stimpacks/env"
	"github.com/readytalk/stim/stimpacks/kubernetes"
	"github.com/readytalk/stim/stimpacks/pagerduty"
	"github.com/readytalk/stim/stimpacks/slack"
	"github.com/readytalk/stim/stimpacks/vault"
	"github.com/readytalk/stim/stimpacks/version"
)

func main() {
	stim := stim.New()
	stim.AddStimpack(discover.New())
	stim.AddStimpack(aws.New())
	stim.AddStimpack(kubernetes.New())
	stim.AddStimpack(pagerduty.New())
	stim.AddStimpack(slack.New())
	stim.AddStimpack(vault.New())
	stim.AddStimpack(version.New())
	stim.AddStimpack(env.New())
	stim.Execute()
}

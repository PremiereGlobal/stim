package main

import (
	"github.com/readytalk/stim/stim"
	// "github.com/readytalk/stim/stimpaks/discover"
	"github.com/readytalk/stim/stimpaks/kubernetes"
	"github.com/readytalk/stim/stimpaks/pagerduty"
	"github.com/readytalk/stim/stimpaks/vault"
	"github.com/readytalk/stim/stimpaks/version"
)

func main() {
	stim := stim.New()
	// stim.AddStimpak(discover.New())
	stim.AddStimpak(kubernetes.New())
	stim.AddStimpak(pagerduty.New())
	stim.AddStimpak(vault.New())
	stim.AddStimpak(version.New())
	stim.Execute()
}

package main

import (
	"github.com/readytalk/stim/stim"
	"github.com/readytalk/stim/stimpaks/pagerduty"
	"github.com/readytalk/stim/stimpaks/vault"
)

func main() {
	stim := stim.New()
	stim.AddStimpak(pagerduty.New())
	stim.AddStimpak(vault.New())
	stim.Execute()
}

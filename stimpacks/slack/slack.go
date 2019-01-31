package slack

import (
	"github.com/readytalk/stim/stim"
)

type Slack struct {
	name string
	stim *stim.Stim
}

func New() *Slack {
	slack := &Slack{name: "slack"}
	return slack
}

func (s *Slack) Name() string {
	return s.name
}

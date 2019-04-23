package slack

import (
	"github.com/PremiereGlobal/stim/stim"
)

const (
	DEFAULT_MESSAGE_USERNAME = "stim"
	DEFAULT_MESSAGE_ICON_URL = "https://vignette.wikia.nocookie.net/fallout/images/7/7e/FoS_stimpak.png/revision/latest"
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

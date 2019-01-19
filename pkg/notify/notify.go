package notify

import (
	VaultApi "github.com/hashicorp/vault/api"
	"github.com/sirupsen/logrus"
)

type Notify struct {
	Options *Options
	log     *logrus.Logger
	vault   *VaultApi.Client
	apikey  string
}

type Options struct {
	Type                 string `mapstructure:"notify-type"`
	PagerdutyServiceName string `mapstructure:"notify-pagerduty"`
	Message              string `mapstructure:"notify-message"`
	Status               string `mapstructure:"notify-status"`
}

func New() *Notify {

	o := &Options{}
	n := &Notify{Options: o}

	return n
}

func (n *Notify) Notify() {

	// Ensure status is set
	if n.Options.Status == "" {
		n.log.Fatal("\"--status -s \" parameter required")
	}

	// Ensure status has a valid value
	if n.Options.Status != "success" && n.Options.Status != "failure" {
		n.log.Fatal("Invalid status \"", n.Options.Status, "\"")
	}

	switch n.Options.Type {
	case "generic":
		n.log.Debug("Notification Type: Generic")

		// Ensure we have a message to send
		if n.Options.Message == "" {
			n.log.Fatal("\"--message\" must be set if using --type \"generic\"")
		}
	case "jenkins-job":
	default:
		n.log.Error("Notification type \"", n.Options.Type, "\" unknown")
	}

}

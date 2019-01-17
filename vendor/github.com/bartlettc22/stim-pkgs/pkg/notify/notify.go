package notify

import (
  "github.com/caarlos0/env"
  log "github.com/sirupsen/logrus"
  // "github.com/bartlettc22/stim-pkgs/pkg/pagerduty"
)

type Notify struct {
  Options *Options
  log *log.Logger
}

type Options struct {
  Type string `mapstructure:"notify-type"`
  // PagerdutyServiceName string `mapstructure:"notify-pagerduty"`
}

type JenkinsJobNotification struct {
  Name string `env:"JOB_NAME" envDefault:"Unknown Job"`
}

func New() *Notify {

  o := &Options{}
  n := &Notify{Options:o}

  return n
}

func (n *Notify) Notify() {

  switch n.Options.Type {
    case "generic":
      n.log.Debug("Notification Type: Generic")
    case "jenkins-job":
      n.log.Debug("Notification Type: Jenkins Job")
      note := JenkinsJobNotification{}
    	err := env.Parse(&note)
    	if err != nil {
    		n.log.Error(err)
    	}
  	default:
  		n.log.Error("Notification type \"", n.Options.Type, "\" unknown")
  }

  // if n.Options.PagerdutyServiceName != "" {
  //
  //   n.log.Debug("Sending Pagerduty Notification to ", n.Options.PagerdutyServiceName)
  //
  //   var event pagerduty.Event
  //   event.ServiceName = n.Options.PagerdutyServiceName
  //
  //   p := pagerduty.New()
  //   response, err := p.CreateEvent(&event)
  //   if err != nil {
  //     log.Error(err)
  //   }
  //
  //   log.Debug(response)
  // }

}

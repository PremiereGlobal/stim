package aws

import (
	"github.com/readytalk/stim/pkg/stimlog"
	// "github.com/aws/aws-sdk-go/aws"
	// "github.com/aws/aws-sdk-go/aws/awserr"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/aws/credentials"
	// "github.com/aws/aws-sdk-go/service/sts"
)

// Aws is the main object
type Aws struct {
	// client *slack.Client
	config *Config
	log    *stimlog.StimLogger
}

type Config struct {
	Token string
}

// New builds a client from the provided config
func New(config *Config, sl *stimlog.StimLogger) (*Aws, error) {

	// client := slack.New(config.Token)

	s := &Aws{config: config}
	if sl != nil {
		s.log = sl
	} else {
		s.log = stimlog.GetLogger()
	}
	return s, nil
}

func (a *Aws) GetCredentials() {
	// a.log.Debug()
}

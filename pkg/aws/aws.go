package aws

import (
	"github.com/readytalk/stim/pkg/log"
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
	log    log.Logger
}

type Config struct {
	Token string
}

// New builds a client from the provided config
func New(config *Config) (*Aws, error) {

	// client := slack.New(config.Token)

	s := &Aws{config: config}

	return s, nil
}

func (a *Aws) GetCredentials() {
	// a.log.Debug()
}

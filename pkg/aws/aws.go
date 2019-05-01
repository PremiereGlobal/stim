package aws

import (
	// 	"github.com/aws/aws-sdk-go/aws"
	// 	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Aws is the main object
type Aws struct {
	config  *Config
	session *session.Session
	log     Logger
}

type Config struct {
	AccessKey string
	SecretKey string
	Log       Logger
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Fatal(args ...interface{})
}

// New builds a client from the provided config
func New(config *Config) (*Aws, error) {

	// Create a new instance of our class
	a := &Aws{config: config, log: config.Log}

	// If credentials were provided, create a new session
	if config.AccessKey != "" && config.SecretKey != "" {
		a.CreateSession(config.AccessKey, config.SecretKey)
	}

	return a, nil
}

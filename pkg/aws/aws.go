package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Aws is the main object
type Aws struct {
	config  *Config
	session *session.Session
}

type Config struct {
	AccessKey string
	SecretKey string
	Logger
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Fatal(args ...interface{})
}

// New builds a client from the provided config
func New(config *Config) (*Aws, error) {

	// Create a new instance of our class
	a := &Aws{config: config}

	// Create a new session based on static IAM credentials that were passed in
	awsCreds := credentials.NewStaticCredentials(config.AccessKey, config.SecretKey, "")
	session, err := session.NewSession(&aws.Config{Credentials: awsCreds})
	if err != nil {
		a.config.Fatal("Error creating AWS session: ", err)
	}

	a.session = session

	// Ensure the credentials are active before we move on
	// Not sure if this should stay here or be optional?
	a.WaitForActiveCreds()

	return a, nil
}

package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func (a *Aws) CreateSession(accessKey string, secretKey string) error {
	awsCreds := credentials.NewStaticCredentials(accessKey, secretKey, "")
	session, err := session.NewSession(&aws.Config{Credentials: awsCreds})
	if err != nil {
		return err
	}
	a.session = session

	return nil
}

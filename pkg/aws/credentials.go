package aws

import (
	"errors"
	"time"

	"github.com/PremiereGlobal/stim/pkg/utils"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sts"
)

// GetFederationToken takes in a name and returns a set of STS Credentials
// based on the current session
func (a *Aws) GetFederationToken(name string, duration time.Duration) *sts.Credentials {

	// Start a new STS session
	s := sts.New(a.session)

	// We are applying admin permissions here which will effectively grant the STS user
	// the same access as the underlying IAM user that is provisioning it.
	// AWS allows the federated user's request only when both the federated user and
	// the IAM user are explicitly allowed to perform the requested action.
	stsUserPolicy := "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Action\":\"*\",\"Resource\":\"*\"}]}"

	// Get and the Federated credentials
	durationSeconds := int64(duration.Seconds())
	output, err := s.GetFederationToken(&sts.GetFederationTokenInput{Name: &name, Policy: &stsUserPolicy, DurationSeconds: &durationSeconds})
	if err != nil {
		a.log.Fatal("Error getting Federation Token: ", err)
	}
	return output.Credentials
}

// WaitForActiveCreds waits for the current session to become valid
// This is useful when IAM credentials were just provisioned and we need to wait
// until they're active to take the next step.
func (a *Aws) WaitForActiveCreds() {

	retryInterval := time.Second * 2
	retryLimit := 20
	successesRequired := 3
	successes := 0

	// Start a new STS session
	s := sts.New(a.session)

	// Here we retry a call to GetCallerIdentity which will return an
	// InvalidClientTokenId error code until the credentials become active
	err := utils.Retry(retryLimit, retryInterval, func() error {

		_, err := s.GetCallerIdentity(&sts.GetCallerIdentityInput{})
		if awserr, ok := err.(awserr.Error); ok {
			if awserr.Code() == "InvalidClientTokenId" {
				a.log.Info("AWS credentials not yet active, waiting...")
				successes = 0
				return err
			} else {
				a.log.Fatal("Error validating AWS credentials: ", err)
			}
		}

		successes += 1
		a.log.Debug("Successful validation check {} of {} reached", successes, successesRequired)
		if successes < successesRequired {
			return errors.New("Haven't reached the required number of consecutive success yet")
		}

		a.log.Info("AWS credentials are active")
		return nil
	})

	// If we've reached this point, the credentials did not become active within
	// the retry limit
	if err != nil {
		a.log.Fatal("Error validating AWS credentials (not active within "+string(retryLimit)+" seconds) ", err)
	}
}

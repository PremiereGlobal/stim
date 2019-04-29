package aws

import (
	"errors"
	"fmt"
	"time"

	stimaws "github.com/PremiereGlobal/stim/pkg/aws"
	"github.com/skratchdot/open-golang/open"
)

// TODO: Move this to a global config
var stimURL = "https://github.com/PremiereGlobal/stim"

// Profile This is our custom profile we'll use to keep track of things
type Profile struct {
	AccessKeyID     string `ini:"aws_access_key_id"`
	SecretAccessKey string `ini:"aws_secret_access_key"`
	SessionToken    string `ini:"aws_session_token"`
	LeaseID         string `ini:"vault_lease_id"`
}

// Login gets IAM or STS credentials
func (a *Aws) Login() error {

	// Create an unauthenticated Aws instance
	a.aws = a.stim.Aws("", "")

	// Create a Vault instance
	a.vault = a.stim.Vault()

	// Prompt the user (or get from arguments) the account and role
	account, role, err := a.GetCredentials()
	if err != nil {
		return err
	}
	a.log.Debug("Account: ", account, " Role: ", role)

	// Looked for a saved profile if desired
	useProfiles := a.stim.GetConfigBool("aws.use-profiles")
	profileName := account + "/" + role
	if useProfiles {
		a.log.Debug("Using AWS profiles")

		// Check if we have a saved profile for the right account/role
		profile := Profile{}
		err := a.aws.MapProfile(profileName, &profile)
		if err != nil {
			return err
		}

		// Validate we got back the expected fields.  If these are missing, the
		// profile is missing or invalid so we'll just generate a new one
		if profile.AccessKeyID != "" && profile.SecretAccessKey != "" && profile.LeaseID != "" {
			a.log.Debug("Profile " + profileName + " found, validating...")
		} else {
			a.log.Debug("Profile " + profileName + " not found or is not familiar...")
		}
	}

	envSource := a.stim.GetConfigBool("env-source")
	stsLogin := a.stim.GetConfigBool("aws-web")
	onlyOutput := a.stim.GetConfigBool("aws-output")

	if stsLogin && a.stim.IsAutomated() {
		a.log.Fatal(errors.New("IsAutomated is detected: web login can not be used."))
	}

	secret, err := a.vault.AWScredentials(account, role)
	if err != nil {
		return err
	}

	accessKey := secret.Data["access_key"].(string)
	secretKey := secret.Data["secret_key"].(string)
	leaseID := secret.LeaseID
	leaseDuration := time.Duration(secret.LeaseDuration) * time.Second
	a.log.Debug("AWS IAM Access Key: " + accessKey)
	a.log.Debug("AWS IAM Access Expiration: " + leaseDuration.String() + " from now")
	a.log.Debug("AWS IAM Vault Lease Id: " + leaseID)

	if useProfiles {

		// Construct our profile
		profile := Profile{
			AccessKeyID:     accessKey,
			SecretAccessKey: secretKey,
			LeaseID:         secret.LeaseID,
		}

		defaultProfile := a.stim.GetConfigBool("aws.default-profile")
		if defaultProfile {
			a.log.Debug("Setting " + profileName + " credentials as default")
		}
		a.aws.SaveProfile(profileName, &profile, defaultProfile)
	}

	if stsLogin {
		aws := a.stim.Aws(accessKey, secretKey)
		federationCreds := aws.GetFederationToken("stim-user")
		a.log.Debug("AWS Federated Access Key: " + *federationCreds.AccessKeyId)
		a.log.Debug("AWS Federated Access Expires: " + federationCreds.Expiration.Sub(time.Now()).String() + " from now")
		loginURL, err := stimaws.CreateAWSLoginURL(*federationCreds.AccessKeyId, *federationCreds.SecretAccessKey, *federationCreds.SessionToken, stimURL)
		a.log.Trace("AWS Console Login URL: " + loginURL)
		if err != nil {
			return err
		}

		if onlyOutput {
			fmt.Print("AWS Console Login URL:\n")
			fmt.Printf("%v\n", loginURL)
		} else {
			err = open.Run(loginURL)
			if err != nil {
				return err
			}
		}
	} else {
		if envSource { // Used for setting AWS credentials in the current environment
			fmt.Println("export AWS_ACCESS_KEY_ID=" + accessKey)
			fmt.Println("export AWS_SECRET_ACCESS_KEY=" + secretKey)
		} else if !useProfiles {
			fmt.Println("AWS_ACCESS_KEY_ID=" + accessKey)
			fmt.Println("AWS_SECRET_ACCESS_KEY=" + secretKey)
		}
	}

	return nil
}

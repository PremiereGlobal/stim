package aws

import (
	"errors"
	"fmt"
	"time"

	awspkg "github.com/PremiereGlobal/stim/pkg/aws"
	"github.com/skratchdot/open-golang/open"
)

// TODO: Move this to a global config
var stimURL = "https://github.com/PremiereGlobal/stim"

// stimProfile This is our custom profile we'll use to keep track of additional
// fields we put in the AWS profile config
type stimProfile struct {
	SessionToken string `ini:"aws_session_token"`
	LeaseID      string `ini:"vault_lease_id"`
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
	a.log.Debug("Account: {} Role: {}", account, role)

	// Looked for a saved profile if desired
	useProfiles := a.stim.ConfigGetBool("aws.use-profiles")
	profileName := account + "/" + role
	if useProfiles {
		a.log.Debug("Using AWS profiles")

		// Check if we have a saved profile for the right account/role
		profile := awspkg.Profile{}
		stimProfile := stimProfile{}
		err := a.aws.MapProfile(profileName, &profile, &stimProfile)
		if err != nil {
			return err
		}

		// Validate we got back the expected fields.  If these are missing, the
		// profile is missing or invalid so we'll just generate a new one
		if profile.AccessKeyID != "" && profile.SecretAccessKey != "" && stimProfile.LeaseID != "" {
			a.log.Debug("Profile {} found", profileName)
		} else {
			a.log.Debug("Profile {} not found or is not familiar", profileName)
		}
	}

	envSource := a.stim.ConfigGetBool("env-source")
	stsLogin := a.stim.ConfigGetBool("aws-web")
	onlyOutput := a.stim.ConfigGetBool("aws-output")

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
	a.log.Debug("AWS IAM Access Key: " + accessKey)
	a.log.Debug("AWS IAM Vault Lease Id: " + leaseID)

	// Get the desired ttl
	ttl, err := time.ParseDuration(a.stim.ConfigGetString("aws.ttl"))
	if err != nil {
		return errors.New(fmt.Sprintf("Error parsing config value aws.ttl: %s", a.stim.ConfigGetString("aws.ttl")))
	}

	// Renew our lease for the requested time
	leaseSecret, err := a.vault.RenewLease(leaseID, ttl)
	if err != nil {
		return err
	}
	a.log.Debug("AWS IAM Access Expiration: " + leaseSecret.String() + " from now")

	if useProfiles {

		// Construct our new base profile
		profile := awspkg.Profile{
			AccessKeyID:     accessKey,
			SecretAccessKey: secretKey,
		}

		// Construct our new stim profile
		stimProfile := stimProfile{
			LeaseID: secret.LeaseID,
		}

		defaultProfile := a.stim.ConfigGetBool("aws.default-profile")
		if defaultProfile {
			a.log.Debug("Setting {} credentials as default", profileName)
		}
		a.aws.SaveProfile(profileName, &profile, defaultProfile, &stimProfile)
	}

	if stsLogin {

		// Get the desired web-ttl
		webTtl, err := time.ParseDuration(a.stim.ConfigGetString("aws.web-ttl"))
		if err != nil {
			return errors.New(fmt.Sprintf("Error parsing config value aws.web-ttl: %s", a.stim.ConfigGetString("aws.web-ttl")))
		}

		a.aws.CreateSession(accessKey, secretKey)
		a.aws.WaitForActiveCreds()
		federationCreds := a.aws.GetFederationToken("stim-user", webTtl)
		a.log.Debug("AWS Federated Access Key: " + *federationCreds.AccessKeyId)
		a.log.Debug("AWS Federated Access Expires: " + federationCreds.Expiration.Sub(time.Now()).String() + " from now")
		loginURL, err := awspkg.CreateAWSLoginURL(*federationCreds.AccessKeyId, *federationCreds.SecretAccessKey, *federationCreds.SessionToken, stimURL)
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

package aws

import (
	"github.com/skratchdot/open-golang/open"

	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TODO: Move this to a global config
var stimURL = "https://github.com/ReadyTalk/stim"

// Login gets IAM or STS credentials
// TODO: Update ~/.aws/ config files for cross shell aws acess
func (a *Aws) Login() error {
	// Create a Vault instance
	a.vault = a.stim.Vault()

	account, role, err := a.GetCredentials()
	if err != nil {
		return err
	}
	a.log.Debug("Account: ", account, " Role: ", role)

	envSource := a.stim.GetConfigBool("env-source")
	stsLogin := a.stim.GetConfigBool("aws-web")
	onlyOutput := a.stim.GetConfigBool("aws-output")
	if stsLogin && a.stim.IsAutomated() {
		a.stim.Fatal(errors.New("IsAutomated is detected: web login can not be used."))
	}

	secret, err := a.vault.AWScredentials(account, role)
	if err != nil {
		return err
	}

	accessKey := secret.Data["access_key"].(string)
	secretKey := secret.Data["secret_key"].(string)
	leaseDuration := time.Duration(secret.LeaseDuration) * time.Second
	a.stim.Debug("AWS IAM Access Key: " + accessKey)
	a.stim.Debug("AWS IAM Access Expiration: " + leaseDuration.String() + " from now")
	a.stim.Debug("AWS IAM Vault Lease Id: " + secret.LeaseID)

	if stsLogin {
		aws := a.stim.Aws(accessKey, secretKey)
		federationCreds := aws.GetFederationToken("stim-user")
		a.stim.Debug("AWS Federated Access Key: " + *federationCreds.AccessKeyId)
		a.stim.Debug("AWS Federated Access Expires: " + federationCreds.Expiration.Sub(time.Now()).String() + " from now")
		loginURL, err := createAWSLoginURL(*federationCreds.AccessKeyId, *federationCreds.SecretAccessKey, *federationCreds.SessionToken)
		// a.stim.Debug("AWS Console Login URL: " + loginURL)
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
		} else {
			fmt.Println("AWS_ACCESS_KEY_ID=" + accessKey)
			fmt.Println("AWS_SECRET_ACCESS_KEY=" + secretKey)
		}
	}

	return nil
}

// createAWSLoginURL returns a federation AWS URL used for wev console login
// This uses AWS Security Token Service (AWS STS) AssumeRole
// More info at: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_enable-console-custom-url.html
// Thanks to Lachlan Donald for the following code: https://github.com/99designs/aws-vault
func createAWSLoginURL(sessionId string, sessionKey string, sessionToken string) (string, error) {
	region := ""
	path := ""
	loginURLPrefix, destination := createRegionalURL(region, path)

	req, err := http.NewRequest("GET", loginURLPrefix, nil)
	if err != nil {
		return "", err
	}

	// Note: This AWS API doesn't validate given info
	jsonBytes, err := json.Marshal(map[string]string{
		"sessionId":    sessionId,
		"sessionKey":   sessionKey,
		"sessionToken": sessionToken,
	})
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("Action", "getSigninToken")
	q.Add("Session", string(jsonBytes))

	req.URL.RawQuery = q.Encode()

	// Note: You can still get a token if you have the wrong credentials
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.New("Failed to create federated token: " + err.Error())
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("Call to getSigninToken failed with " + resp.Status)
	}

	var respParsed map[string]string

	if err = json.Unmarshal([]byte(body), &respParsed); err != nil {
		return "", errors.New("Failed to parse response from getSigninToken: " + err.Error())
	}

	signinToken, ok := respParsed["SigninToken"]
	if !ok {
		return "", errors.New("Expected a response with SigninToken")
	}

	loginURL := fmt.Sprintf(
		"%s?Action=login&Issuer=%s&Destination=%s&SigninToken=%s",
		loginURLPrefix,
		url.QueryEscape(stimURL),
		url.QueryEscape(destination),
		url.QueryEscape(signinToken),
	)

	return loginURL, nil
}

// createRegionalURL create the needed regional AWS URL
func createRegionalURL(region string, path string) (string, string) {
	loginURLPrefix := "https://signin.aws.amazon.com/federation"
	destination := "https://console.aws.amazon.com/"

	if region != "" {
		destinationDomain := "console.aws.amazon.com"
		switch {
		case strings.HasPrefix(region, "cn-"):
			loginURLPrefix = "https://signin.amazonaws.cn/federation"
			destinationDomain = "console.amazonaws.cn"
		case strings.HasPrefix(region, "us-gov-"):
			loginURLPrefix = "https://signin.amazonaws-us-gov.com/federation"
			destinationDomain = "console.amazonaws-us-gov.com"
		}
		if path != "" {
			destination = fmt.Sprintf(
				"https://%s.%s/%s?region=%s",
				region, destinationDomain, path, region,
			)
		} else {
			destination = fmt.Sprintf(
				"https://%s.%s/console/home?region=%s",
				region, destinationDomain, region,
			)
		}
	}
	return loginURLPrefix, destination
}

package aws

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// CreateAWSLoginURL returns a federation AWS URL used for wev console login
// This uses AWS Security Token Service (AWS STS) AssumeRole
// More info at: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_providers_enable-console-custom-url.html
// Thanks to Lachlan Donald for the following code: https://github.com/99designs/aws-vault
func CreateAWSLoginURL(sessionId string, sessionKey string, sessionToken string, issuer string) (string, error) {
	region := ""
	path := ""
	loginURLPrefix, destination := CreateRegionalURL(region, path)

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
		url.QueryEscape(issuer),
		url.QueryEscape(destination),
		url.QueryEscape(signinToken),
	)

	return loginURL, nil
}

// CreateRegionalURL creates signin and console URLs based on the region/path
// provided
func CreateRegionalURL(region string, path string) (string, string) {
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

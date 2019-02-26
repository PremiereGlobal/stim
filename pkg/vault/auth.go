package vault

import (
	"github.com/hashicorp/vault/command/token"
	"golang.org/x/crypto/ssh/terminal"

	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
)

// This is the main Login function
func (v *Vault) Login() error {
	v.tokenHelper = token.InternalTokenHelper{}

	if v.client.Token() == "" { // If no environment token set
		// Reading token from user's dot file
		token, err := v.tokenHelper.Get()
		if err != nil {
			return err
		}
		if token != "" {
			v.client.SetToken(token)
			v.Debug("Reading token from: " + v.tokenHelper.Path())
		} else { // If we still can not find the token
			v.Debug("No token found. Trying to login.")
			err = v.userLogin()
			if err != nil {
				return err
			}
		}
	} else {
		v.Debug("Reading token from environment 'VAULT_TOKEN'")
	}

	// Test token and see if a vault login is needed
	loginToVault := false
	r := v.client.NewRequest("GET", "/v1/auth/token/lookup-self")
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := v.client.RawRequestWithContext(ctx, r) // Access to resp is nice
	if err != nil {
		if resp.StatusCode == 403 {
			v.Debug("Got permission denied. Trying to login.")
			loginToVault = true
		} else {
			return v.parseError(err)
		}
	}
	defer resp.Body.Close()

	if loginToVault == true {
		v.Debug("Need to login to Vault")
		v.userLogin()
	}

	return nil
}

func (v *Vault) isCurrentTokenValid() {

}

func (v *Vault) userLogin() error {
	// Sadly we will assume LDAP login (for now)
	// Maybe someday vault will allow anonymous access to "vault auth list"

	if v.config.Noprompt == true {
		return errors.New("No interactive prompt is set, but user input is required to continue")
	}

	username, password, err := v.getCredentials()
	if err != nil {
		return err
	}

	// No hacking: Test username
	// https://stackoverflow.com/questions/6949667/what-are-the-real-rules-for-linux-usernames-on-centos-6-and-rhel-6
	// match, err := regexp.MatchString("^[a-z_][a-z0-9_]{0,30}$", username)
	// check(err)
	// if match != true {
	// 	// Safe to exit, we know we are in a user interface
	// 	log.Fatal("Username does not match BSD 4.3 standards (32 character string 0f [a-z0-9_])")
	// }

	// Login with LDAP and create a token
	secret, err := v.client.Logical().Write("auth/ldap/login/"+username, map[string]interface{}{
		"password": password,
	})
	if err != nil {
		v.Debug("Do you have a bad username or password?")
		return v.parseError(err)
	}
	v.client.SetToken(secret.Auth.ClientToken)

	// Write token to user's dot file
	err = v.tokenHelper.Store(secret.Auth.ClientToken)
	if err != nil {
		return v.parseError(err)
	}

	// Lookup the token to get the entity ID
	secret, err = v.client.Auth().Token().Lookup(v.client.Token())
	if err != nil {
		return v.parseError(err)
	}
	// spew.Dump(secret)
	// entityID := secret.Data["entity_id"].(string)
	// Debug("Vault entity ID: ", entityID)

	return nil
}

// Gather username and password from the user
// Could also use: github.com/hashicorp/vault/helper/password
func (v *Vault) getCredentials() (string, string, error) {
	fmt.Println("Vault needs your LDAP Linux user/pass.")
	if v.config.Username != "" {
		fmt.Printf("Username (%s): ", v.config.Username)
	} else {
		fmt.Printf("Username: ")
	}
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	if len(username) <= 0 { // If user just clicked enter
		if v.config.Username == "" { // If there also isn't default
			return "", "", v.newError("No username given")
		}
		username = v.config.Username
	} else {
		v.config.Username = username
	}

	fmt.Print("Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", v.parseError(err)
	}
	fmt.Println("")
	password := string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

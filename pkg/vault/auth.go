package vault

import (
	"github.com/hashicorp/vault/command/token"
	"github.com/readytalk/stim/pkg/log"
	"golang.org/x/crypto/ssh/terminal"
	// "github.com/davecgh/go-spew/spew"

	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"syscall"
)

// Login will authenticate the user with Vault
// Will detect if user needs to re-login
func (v *Vault) Login() error {

	// get the token from the user's environment
	v.tokenHelper = token.InternalTokenHelper{}
	if v.client.Token() != "" {
		log.Debug("Reading token from environment 'VAULT_TOKEN'")
	} else { // If no environment token set
		// Reading token from user's dot file
		token, err := v.tokenHelper.Get()
		if err != nil {
			return err
		}
		if token != "" {
			v.client.SetToken(token)
			log.Debug("Reading token from: " + v.tokenHelper.Path())
		} else { // If we still can not find the token
			log.Debug("No token found. Trying to login.")
			err = v.userLogin()
			if err != nil {
				return err
			}
		}
	}

	isTokenValid := v.isCurrentTokenValid()
	if isTokenValid == false {
		log.Debug("Need to login to Vault")
		v.userLogin()
	}

	return nil
}

// isCurrentTokenValid returns flase if user needs to relogin
// I am not happy with this way of testing the token.
// This function doesn't check for errors.
func (v *Vault) isCurrentTokenValid() bool {
	secret, _ := v.client.Auth().Token().LookupSelf()
	if secret == nil {
		return false
	}
	return true
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
		log.Debug("Do you have a bad username or password?")
		return err
	}
	v.client.SetToken(secret.Auth.ClientToken)

	// Write token to user's dot file
	err = v.tokenHelper.Store(secret.Auth.ClientToken)
	if err != nil {
		return err
	}

	// Lookup the token to get the entity ID
	secret, err = v.client.Auth().Token().Lookup(v.client.Token())
	if err != nil {
		return err
	}
	// spew.Dump(secret)
	entityID := secret.Data["entity_id"].(string)
	log.Debug("Vault entity ID: ", entityID)

	v.newLogin = true // Set if we had to prompt user for a login

	return nil
}

// IsNewLogin will help high level funcs know if a login prompt was used
func (v *Vault) IsNewLogin() bool {
	return v.newLogin
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
			return "", "", errors.New("No username given")
		}
		username = v.config.Username
	} else {
		v.config.Username = username
	}

	fmt.Print("Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	fmt.Println("")
	password := string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

// RenewToken will increase the lease duration of users token
func (v *Vault) RenewToken(time string) error {
	log.Debug("Renewing token for '", time, "'")
	_, err := v.client.Logical().Write("auth/token/renew-self", map[string]interface{}{
		"increment": time,
	})
	if err != nil {
		return err
	}

	return nil
}

package vault

import (
	"bufio"
	"errors"
	"github.com/hashicorp/vault/command/token"
	// 	"context"
	"fmt"
	// 	VaultApi "github.com/hashicorp/vault/api"
	// 	VaultToken "github.com/hashicorp/vault/command/token"
	// 	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"os/user"
	// 	"reflect"
	// 	"regexp"
	"strings"
	"syscall"
	// 	"time"
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
			// log.Debug("Reading token from user's dot file")
		} else { // If we still can not find the token
			// log.Debug("No token found. Trying to login.")
			err = v.userLogin()
			if err != nil {
				return err
			}
		}
	}
	// else {
	// 	log.Debug("Reading token from user's environment variable")
	// }

	// Test token and see if a vault login is needed
	// loginToVault := false
	// r := v.client.NewRequest("GET", "/v1/auth/token/lookup-self")
	// ctx, cancelFunc := context.WithCancel(context.Background())
	// defer cancelFunc()
	// resp, err := v.client.RawRequestWithContext(ctx, r) // Access to resp is nice
	// if err != nil {
	//   if resp.StatusCode == 403 {
	//     log.Debug("Got permission denied. Trying to login.")
	//     loginToVault = true
	//   } else {
	//     log.Error(err)
	//   }
	// }
	// defer resp.Body.Close()
	//
	// if loginToVault == true {
	//   log.Debug("Need to login to Vault")
	//   v.userLoginPrompt()
	// }

	return nil
}

func (v *Vault) isCurrentTokenValid() {

}

func (v *Vault) userLogin() error {
	// Sadly we will for now assume LDAP login
	// Maybe someday vault will allow anonymous access to "vault auth list"

	if v.config.Noprompt == true {
		return errors.New("No interactive prompt is set, but user input is required to continue")
		// log.Error("No interactive prompt is set, but user input is required to continue")
		// os.Exit(1)
	}

	username, password, err := v.getCredentials()
	if err != nil {
		return err
	}
	// fmt.Printf("Username: %s, Password: %s\n", username, password)

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
		return err
		// log.Info("Do you have a bad username or password?")
		// log.Fatal(err)
	}
	v.client.SetToken(secret.Auth.ClientToken)

	err = v.tokenHelper.Store(secret.Auth.ClientToken)
	if err != nil {
		return err
		// log.Fatal(err)
	}

	// Lookup the token to get the entity ID
	// secret, err = v.client.Auth().Token().Lookup(v.client.Token())
	// check(err)
	// entityID := secret.Data["entity_id"].(string)
	// log.Debug("Vault entity ID: ", entityID)

	return nil
}

func (v *Vault) getCredentials() (string, string, error) {

	user, _ := user.Current()
	// v.stim.DebugError(err)

	fmt.Println("Vault needs your LDAP Linux user/pass.")
	fmt.Printf("Username (%s): ", user.Username)
	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	if len(username) <= 0 {
		username = user.Username
	}

	fmt.Print("Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}
	// check(err)
	password := string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

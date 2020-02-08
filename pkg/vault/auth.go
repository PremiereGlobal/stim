package vault

import (
	"github.com/hashicorp/vault/command/token"
	"golang.org/x/crypto/ssh/terminal"

	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"syscall"
)

// Login will authenticate the user with Vault
// Will detect if user needs to re-login
func (v *Vault) Login() error {
	// get the token from the user's environment
	if v.client.Token() != "" {
		v.log.Debug("Reading token from environment 'VAULT_TOKEN'")
	} else { // If no environment token set
		// Reading token from user's dot file
		v.tokenHelper = token.InternalTokenHelper{}
		token, err := v.tokenHelper.Get()
		if err != nil {
			return v.parseError(err).(error)
		}

		if token != "" {
			v.log.Debug("Reading token from: " + v.tokenHelper.Path())
			v.client.SetToken(token)
		}
	}

	// Check if any existing token is valid
	// If not, prompt for login
	isTokenValid := v.isCurrentTokenValid()

	if isTokenValid == false {
		v.log.Debug("No valid tokens found, need to login")
		err := v.userLogin()
		if err != nil {
			return v.parseError(err)
		}
	}

	return nil
}

// GetToken returns the raw token
func (v *Vault) GetToken() (string, error) {
	if token := v.client.Token(); token != "" {
		return token, nil
	}

	return "", errors.New("No token set")
}

// GetMeta returns the username metadata for the current token
func (v *Vault) GetUsername() (string, error) {

	secret, err := v.client.Auth().Token().LookupSelf()
	if err != nil {
		return "", v.parseError(err)
	}

	metadata, err := secret.TokenMetadata()
	if err != nil {
		return "", v.parseError(err)
	}

	if val, ok := metadata["username"]; ok {
		return val, nil
	}

	return "", errors.New("No username set")
}

// isCurrentTokenValid returns false if user needs to re-login
// I am not happy with this way of testing the token.
func (v *Vault) isCurrentTokenValid() bool {
	duration, err := v.GetCurrentTokenTTL()
	if err != nil || duration <= 0 {
		return false
	}

	v.log.Debug("Current token is valid for {}", duration.String())

	return true
}

// userLogin authenticates with Vault and obtains a user token
func (v *Vault) userLogin() error {

	if v.config.Noprompt == true {
		return errors.New("No interactive prompt is set, but user input is required to continue")
	}

	username, password, err := v.getCredentials()
	if err != nil {
		return err
	}

	// Login and obtain a token
	authPath := path.Join("auth/", v.config.AuthPath, "/login/", username)
	secret, err := v.client.Logical().Write(authPath, map[string]interface{}{
		"password": password,
	})
	if err != nil {
		v.log.Debug("Do you have a bad username or password?")
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
	entityID := secret.Data["entity_id"].(string)
	v.log.Debug("Vault entity ID: ", entityID)

	v.newLogin = true // Set if we had to prompt user for a login

	return nil
}

// IsNewLogin will help high level funcs know if a login prompt was used
func (v *Vault) IsNewLogin() bool {
	return v.newLogin
}

// getCredentials gathers username and password from the user
// Could also use: github.com/hashicorp/vault/helper/password
func (v *Vault) getCredentials() (string, string, error) {

	var username string
	fmt.Println("Please enter your [" + v.config.AuthPath + "] credentials")
	if v.config.UsernameSkipPrompt && v.config.Username != "" {
		v.log.Debug("Skipping username prompt. Using config value 'vault-username:{}'", v.config.Username)
		fmt.Printf("Username: %s\n", v.config.Username)
		username = v.config.Username
	} else {
		if v.config.Username != "" {
			fmt.Printf("Username (%s): ", v.config.Username)
		} else {
			fmt.Printf("Username: ")
		}

		reader := bufio.NewReader(os.Stdin)
		username, _ = reader.ReadString('\n')
		username = strings.TrimSpace(username)

		if len(username) <= 0 { // If user just clicked enter
			if v.config.Username == "" { // If there also isn't default
				return "", "", v.newError("No username given").(error)
			}
			username = v.config.Username
		} else {
			v.config.Username = username
		}
	}

	fmt.Print("Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", v.parseError(err).(error)
	}
	fmt.Println("")
	password := string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

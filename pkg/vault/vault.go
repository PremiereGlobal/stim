package vault

import (
	"bufio"
	"context"
	"fmt"
	VaultApi "github.com/hashicorp/vault/api"
	VaultToken "github.com/hashicorp/vault/command/token"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"os/user"
	"reflect"
	"strings"
	"syscall"
	"time"
  "regexp"
)

type Client struct {
	Config *Config
	client *VaultApi.Client
  tokenHelper VaultToken.InternalTokenHelper
}

type Config struct {
  Noprompt bool
	Address string `vrequired:"true" mapstructure:"vault-address"`
}

var log *logrus.Logger

func check(e error, exit ...bool) { // This helper will streamline our error checks below.
	if e != nil {
		log.Error(e)
    if len(exit) != 0 {
      if exit[0] == true {
        os.Exit(1)
      }
    }
  }
}

func SetLogger(givenLog *logrus.Logger) {
	log = givenLog
}

func NewVault() *Client {
	config := &Config{}

	c := &Client{Config: config}

	return c
}

func (v *Client) Setup() {
	// Make sure all needed variables are set`
	checkRequired(v.Config)

	// Configure new Vault Client
	config := VaultApi.DefaultConfig()
	config.Address = v.Config.Address // Since we read the env we can override
	// config.HttpClient.Timeout = 60    // No need to wait over a minite from default

	var err error
	v.client, err = VaultApi.NewClient(config)
	check(err)
  v.tokenHelper = VaultToken.InternalTokenHelper{}

	// Test the Vault server
	_, err = v.isVaultHealthy()
	check(err, true)

	if v.client.Token() == "" { // If no environment token set
		// Reading token from user's dot file
		token, _ := v.tokenHelper.Get()
    check(err)
		v.client.SetToken(token)
	}

	if v.client.Token() == "" { // If we still can not find the token
		log.Debug("No token found. Trying to login.")
    v.userLoginPrompt()
	} else {
		log.Debug("Reading token from user's environment variable")
	}

	// Test token and see if a vault login is needed
  loginToVault := false
	r := v.client.NewRequest("GET", "/v1/auth/token/lookup-self")
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	resp, err := v.client.RawRequestWithContext(ctx, r) // Access to resp is nice
	if err != nil {
		if resp.StatusCode == 403 {
			log.Debug("Got permission denied. Trying to login.")
      loginToVault = true
		} else {
			log.Error(err)
		}
	}
	defer resp.Body.Close()

  if loginToVault == true {
    log.Debug("Need to login to Vault")
    v.userLoginPrompt()
  }
}

func (v *Client) isVaultHealthy() (bool, error) {
	result, err := v.client.Sys().Health()
	if err != nil {
		return false, err
	}

	log.Debug("Vault server info from (", v.client.Address(), ")")
	log.Debug("  Initialized: ", result.Initialized)
	log.Debug("  Sealed: ", result.Sealed)
	log.Debug("  Standby: ", result.Standby)
	log.Debug("  Version: ", result.Version)
	log.Debug("  ClusterName: ", result.ClusterName)
	log.Debug("  ClusterID: ", result.ClusterID)
	log.Debug("  ServerTime: (", result.ServerTimeUTC, ") ", time.Unix(result.ServerTimeUTC, 0).UTC())
	log.Debug("  Standby: ", result.Standby)

	return true, nil
}

func (v *Client) userLoginPrompt() (bool) {
  // Sadly we will for now assume LDAP login
  // Maybe someday vault will allow anonymous access to "vault auth list"

  if v.Config.Noprompt == true {
    log.Error("No interactive prompt is set, but Vault login is required to continue")
    os.Exit(1)
  }

  username, password := credentials()
  fmt.Printf("Username: %s, Password: %s\n", username, password)

  // No hacking: Test username
  // https://stackoverflow.com/questions/6949667/what-are-the-real-rules-for-linux-usernames-on-centos-6-and-rhel-6
  match, err := regexp.MatchString("^[a-z_][a-z0-9_]{0,30}$", username)
  check(err)
  if match != true {
    // Safe to exit, we know we are in a user interface
    log.Fatal("Username does not match BSD 4.3 standards (32 character string 0f [a-z0-9_])")
  }

  // Login with LDAP and create a token
	secret, err := v.client.Logical().Write("auth/ldap/login/"+username, map[string]interface{}{
		"password": password,
	})
  if err != nil {
    log.Info("Do you have a bad username or password?")
    log.Fatal(err)
  }
  v.client.SetToken(secret.Auth.ClientToken)

  err = v.tokenHelper.Store(secret.Auth.ClientToken)
  if err != nil {
    log.Fatal(err)
  }

	// Lookup the token to get the entity ID
  secret, err = v.client.Auth().Token().Lookup(v.client.Token())
  check(err)
  entityID := secret.Data["entity_id"].(string)
  log.Debug("Vault entity ID: ", entityID)

  return false
}

// Gather username and password from the user
// Could also use: github.com/hashicorp/vault/helper/password
func credentials() (string, string) {
	user, err := user.Current()
	check(err)

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
	check(err)
	password := string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password)
}

func checkRequired(spec *Config) {
	t := reflect.TypeOf(*spec)

	for i := 0; i < t.NumField(); i++ {
		// Get the field, returns https://golang.org/pkg/reflect/#StructField
		field := t.Field(i)

		// Get the field tag value
		tag := field.Tag.Get("vrequired")

		if tag == "true" && field.Type.Name() != "bool" {
			r := reflect.ValueOf(spec)
			fieldValue := reflect.Indirect(r).FieldByName(field.Name)
			if fieldValue.String() == "" {
				log.Fatal(field.Name + " required but not set. Use environment variable " + field.Tag.Get("envconfig") + " or command line options: --" + field.Tag.Get("long") + ", -" + field.Tag.Get("short"))
			}
		}
	}
}

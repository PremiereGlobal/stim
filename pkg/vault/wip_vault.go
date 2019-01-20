package vault

//
// import (
// 	"bufio"
// 	"fmt"
// 	VaultApi "github.com/hashicorp/vault/api"
// 	homedir "github.com/mitchellh/go-homedir"
// 	"github.com/sirupsen/logrus"
// 	"golang.org/x/crypto/ssh/terminal"
// 	"os"
// 	"os/user"
// 	"reflect"
// 	"strings"
// 	"syscall"
// 	"time"
// )
//
// type Client struct {
// 	Config *Config
// 	client *VaultApi.Client
// }
//
// type Config struct {
// 	Address    string `vrequired:"true" mapstructure:"vault-address"`
// 	Token      string `mapstructure:"vault-token"`
// 	TokenFile  string
// 	SkipVerify bool
// }
//
// var log *logrus.Logger
//
// func check(e error) { // This helper will streamline our error checks below.
// 	if e != nil {
// 		log.Error(e)
// 	}
// }
//
// func SetLogger(givenLog *logrus.Logger) {
// 	log = givenLog
// }
//
// func NewVault() *Client {
// 	config := &Config{}
//
// 	home, err := homedir.Dir() // Find home directory.
// 	check(err)
//
// 	config.TokenFile = home + "/.vault-token"
//
// 	c := &Client{Config: config}
//
// 	return c
// }
//
// func (v *Client) Login() {
// 	// Make sure all needed variables are set`
// 	checkRequired(v.Config)
//
// 	// Configure new Vault Client
// 	conf := &VaultApi.Config{Address: v.Config.Address}
// 	tlsConf := &VaultApi.TLSConfig{Insecure: v.Config.SkipVerify}
// 	conf.ConfigureTLS(tlsConf)
// 	client, err := VaultApi.NewClient(conf)
// 	check(err)
// 	v.client = client
//
// 	// Test the Vault server
// 	_, err = v.isVaultHealthy()
// 	check(err)
//
// 	if v.Config.Token == "" {
// 		username, password := credentials()
// 		fmt.Printf("Username: %s, Password: %s\n", username, password)
// 	}
//
// 	client.SetToken(v.Config.Token)
// 	_, err = client.Auth().Token().LookupSelf()
// 	check(err)
//
// 	// Work in progress. Need to complete login func
// }
//
// func (v *Client) isVaultHealthy() (bool, error) {
// 	result, err := v.client.Sys().Health()
// 	if err != nil {
// 		return false, err
// 	}
//
// 	log.Debug("Vault server info from (", v.client.Address(), ")")
// 	log.Debug("  Initialized: ", result.Initialized)
// 	log.Debug("  Sealed: ", result.Sealed)
// 	log.Debug("  Standby: ", result.Standby)
// 	log.Debug("  Version: ", result.Version)
// 	log.Debug("  ClusterName: ", result.ClusterName)
// 	log.Debug("  ClusterID: ", result.ClusterID)
// 	log.Debug("  ServerTime: (", result.ServerTimeUTC, ") ", time.Unix(result.ServerTimeUTC, 0).UTC())
// 	log.Debug("  Standby: ", result.Standby)
//
// 	return true, nil
// }
//
// // Gather username and password from the user
// // Could also use: github.com/hashicorp/vault/helper/password
// func credentials() (string, string) {
// 	user, err := user.Current()
// 	check(err)
//
// 	fmt.Println("Vault needs your LDAP Linux user/pass.")
// 	fmt.Printf("Username (%s): ", user.Username)
// 	reader := bufio.NewReader(os.Stdin)
// 	username, _ := reader.ReadString('\n')
// 	username = strings.TrimSpace(username)
// 	if len(username) <= 0 {
// 		username = user.Username
// 	}
//
// 	fmt.Print("Password: ")
// 	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
// 	check(err)
// 	password := string(bytePassword)
//
// 	return strings.TrimSpace(username), strings.TrimSpace(password)
// }
//
// func checkRequired(spec *Config) {
// 	t := reflect.TypeOf(*spec)
//
// 	for i := 0; i < t.NumField(); i++ {
// 		// Get the field, returns https://golang.org/pkg/reflect/#StructField
// 		field := t.Field(i)
//
// 		// Get the field tag value
// 		tag := field.Tag.Get("vrequired")
//
// 		if tag == "true" && field.Type.Name() != "bool" {
// 			r := reflect.ValueOf(spec)
// 			fieldValue := reflect.Indirect(r).FieldByName(field.Name)
// 			if fieldValue.String() == "" {
// 				log.Fatal(field.Name + " required but not set. Use environment variable " + field.Tag.Get("envconfig") + " or command line options: --" + field.Tag.Get("long") + ", -" + field.Tag.Get("short"))
// 			}
// 		}
// 	}
// }

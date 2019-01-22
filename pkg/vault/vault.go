package vault

import (
	"errors"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/token"
	// "log"
	"fmt"
)

type Vault struct {
	client      *api.Client
	tokenHelper token.InternalTokenHelper
	config      *Config
}

type Config struct {
	Noprompt bool
	Address  string //`vrequired:"true" mapstructure:"vault-address"`
	Logger
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
}

func (v *Vault) Debug(message string) {
	if v.config.Logger != nil {
		v.config.Debug(message)
	}
}

func (v *Vault) Info(message string) {
	if v.config.Logger != nil {
		v.config.Info(message)
	} else {
		fmt.Println(message)
	}
}

func New(config *Config) (*Vault, error) {
	// Ensure that the Vault address is set
	if config.Address == "" {
		return nil, errors.New("Vault address not set")
	}

	v := &Vault{config: config}

	// Configure new Vault Client
	apiConfig := api.DefaultConfig()
	apiConfig.Address = v.config.Address // Since we read the env we can override
	// config.HttpClient.Timeout = 60    // No need to wait over a minite from default

	// Create our new API client
	var err error
	v.client, err = api.NewClient(apiConfig)
	if err != nil {
		return nil, err
	}

	// Ensure Vault is up and Healthy
	_, err = v.isVaultHealthy()
	if err != nil {
		return nil, err
	}

	// Run Login logic
	err = v.Login()
	if err != nil {
		return nil, err
	}

	return v, nil
}

// )

// type Client struct {
// 	Config      *Config
// 	client      *VaultApi.Client
// 	tokenHelper VaultToken.InternalTokenHelper
// }
//

// func (v *Vault) InitClient() error {
// 	// Initialize client
// 	client, err := api.NewClient(nil)
// 	if err != nil {
// 		return err
// 	}
// 	tokenHelper := token.InternalTokenHelper{}
// 	token, err := tokenHelper.Get()
// 	if err != nil {
// 		return err
// 	}
// 	client.SetToken(token)
// 	v.client = client
//
// 	return nil
// }
//
// func (v *Vault) Login() {
//
// }

// =======
// func check(e error, exit ...bool) { // This helper will streamline our error checks below.
// 	if e != nil {
// 		log.Error(e)
// 		if len(exit) != 0 {
// 			if exit[0] == true {
// 				os.Exit(1)
// 			}
// 		}
// 	}
// }

// func SetLogger(givenLog *logrus.Logger) {
// 	log = givenLog
// }

// func NewVault() *Client {
// 	config := &Config{}
//
// 	c := &Client{Config: config}
//
// 	return c
// }

// func (v *Client) Setup() {
//
// }

// Gather username and password from the user
// Could also use: github.com/hashicorp/vault/helper/password

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

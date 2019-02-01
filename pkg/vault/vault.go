package vault

import (
	"errors"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/token"

	"fmt"
)

type Vault struct {
	client      *api.Client
	tokenHelper token.InternalTokenHelper
	config      *Config
}

type Config struct {
	Noprompt bool
	Address  string
	Username string
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

	// v.Debug("Username: "+ v.config.username)

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

func (v *Vault) GetUser() string {
	return v.config.Username
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

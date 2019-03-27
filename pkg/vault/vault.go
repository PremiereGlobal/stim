package vault

import (
	"fmt"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/token"
	"github.com/readytalk/stim/pkg/log"
	"errors"
	"time"
)

type Vault struct {
	client      *api.Client
	config      *Config
	tokenHelper token.InternalTokenHelper
	newLogin    bool
}

type Config struct {
	Noprompt bool
	Address  string
	Username string
	Timeout  int
  InitialTokenDuration time.Duration
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

	v := &Vault{config: config}
	log.SetLogger(config.Log)

	// Ensure that the Vault address is set
	if config.Address == "" {
		return nil, v.newError("Vault address not set")
	}

	// Configure new Vault Client
	apiConfig := api.DefaultConfig()
	apiConfig.Address = v.config.Address // Since we read the env we can override
	apiConfig.Timeout = time.Duration(v.config.Timeout) * time.Second

	// Create our new API client
	var err error
	v.client, err = api.NewClient(apiConfig)
	if err != nil {
		return nil, v.parseError(err)
	}

	// Ensure Vault is up and Healthy
	_, err = v.isVaultHealthy()
	if err != nil {
		return nil, v.parseError(err)
	}

	// Run Login logic
	err = v.Login()
	if err != nil {
		return nil, v.parseError(err)
	}

	// If user wants, extend the token timeout
	if v.IsNewLogin() {
		if v.config.InitialTokenDuration > 0 {
			log.Debug("Token duration set to: ", v.config.InitialTokenDuration)
			_, err = v.client.Auth().Token().RenewSelf(int(v.config.InitialTokenDuration))
			if err != nil {
				return nil, err
			}
		}
	}

	return v, nil
}

func (v *Vault) GetUser() string {
	return v.config.Username
}

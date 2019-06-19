package vault

import (
	"time"

	"github.com/PremiereGlobal/stim/pkg/stimlog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/command/token"
)

type Vault struct {
	client      *api.Client
	config      *Config
	tokenHelper token.InternalTokenHelper
	newLogin    bool
	log         Logger
}

type Config struct {
	AuthPath             string
	Noprompt             bool
	Address              string
	Username             string
	Timeout              time.Duration
	InitialTokenDuration time.Duration
	Log                  Logger
}

type Logger interface {
	Debug(...interface{})
	Warn(...interface{})
	Fatal(...interface{})
}

func New(config *Config) (*Vault, error) {
	v := &Vault{config: config}
	if config.Log != nil {
		v.log = config.Log
	} else {
		v.log = stimlog.GetLogger()
	}

	// Ensure that the Vault address is set
	if config.Address == "" {
		return nil, v.newError("Vault address not set")
	}

	// Set a default auth method to ldap if not set
	// TODO: change this to token auth as it is more generic
	if config.AuthPath == "" {
		config.AuthPath = "ldap"
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
		return nil, err
	}

	// Run Login logic
	err = v.Login()
	if err != nil {
		return nil, err
	}

	// If user wants, extend the token timeout
	if v.IsNewLogin() {
		if v.config.InitialTokenDuration > 0 {
			v.log.Debug("Attempting to set token duration to ", v.config.InitialTokenDuration)
			_, err = v.client.Auth().Token().RenewSelf(int(v.config.InitialTokenDuration.Seconds()))
			if err != nil {
				return nil, v.parseError(err)
			}

			// Show the actual TTL and warn if different from requested
			actualDuration, err := v.GetCurrentTokenTTL()
			if err != nil {
				return nil, err
			}
			v.log.Debug("Current token is valid for {}", actualDuration.String())

			// If duration is not within threshold of the requested duration, warn the user
			if actualDuration.Seconds() <= (v.config.InitialTokenDuration.Seconds() - 15) {
				v.log.Warn("Unable to set token duration to {}.  Actual value: {}", v.config.InitialTokenDuration.String(), actualDuration.String())
			}
		}
	}

	return v, nil
}

func (v *Vault) GetUser() string {
	return v.config.Username
}

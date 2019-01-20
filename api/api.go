package api

import (
	"github.com/readytalk/stim/pkg/pagerduty"
	// "github.com/readytalk/stim/pkg/prometheus"
	"github.com/readytalk/stim/pkg/vault"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Api struct {
	Config *viper.Viper
	Log    Log
}

type Log logrus.FieldLogger

func New(viper *viper.Viper) *Api {
	api := &Api{Config: viper}
	return api
}

func (a *Api) BindLogger(log Log) {
	a.Log = log
}

func (a *Api) Pagerduty() *pagerduty.Pagerduty {
	a.Log.Debug("API: Creating Pagerduty")
	vaultPath := a.Config.Get("pagerduty.vault-apikey-path").(string)
	vaultKey := a.Config.Get("pagerduty.vault-apikey-key").(string)
	a.Log.Debug("API: Fetching Pagerduty API key from Vault `", vaultPath, "`")
	vault := a.Vault()
	apikey, err := vault.GetSecretKey(vaultPath, vaultKey)
	if err != nil {
		a.Log.Fatal("API Pagerduty: Error getting API key from Vault: ", err)
	}
	pagerduty := pagerduty.New(apikey)
	return pagerduty
}

//
// func (a *Api) Prometheus() *prometheus.Prometheus {
// 	prometheus := prometheus.New()
// 	return prometheus
// }

func (a *Api) Vault() *vault.Vault {
	vault := vault.New()
	err := vault.InitClient()
	if err != nil {
		a.Log.Fatal("API Vault: Error Initializaing Client: ", err)
	}
	return vault
}

func (a *Api) BindCommand(command *cobra.Command, parentCommand *cobra.Command) {
	parentCommand.AddCommand(command)
}

// func (v *Vault) AddFlags() {
//   // Get the Vault Address from the environment and command line
//   vaultCmd.PersistentFlags().StringP("address", "a", "http://127.0.0.1:8200", "Vault URL")
//   v.api.Config.BindEnv("vault-address", "VAULT_ADDR")
//   v.api.Config.BindPFlag("vault-address", vaultCmd.PersistentFlags().Lookup("address"))
// }

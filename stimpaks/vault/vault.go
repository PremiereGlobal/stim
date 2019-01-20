package vault

import (
	"github.com/readytalk/stim/api"
	"github.com/spf13/cobra"
)

type Vault struct {
	api *api.Api
	log api.Log
}

func New(api *api.Api) *Vault {
	vault := &Vault{api: api, log: api.Log}
	return vault
}

func (v *Vault) Bind(rootCmd *cobra.Command) {
	// v.BindLogger(log)
	var vaultCmd = &cobra.Command{
		Use:   "vault",
		Short: "Vault helper",
		Long:  "Vault LDAP login and AWS access",
		// PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// 	vault.SetLogger(log) // Setup logging
		//
		// 	vaultHandler = vault.NewVault() // Setup a vault object
		//
		// 	// Get the Vault token only from the environment variable and ~/.vault-token file
		// 	if _, ok := os.LookupEnv("VAULT_TOKEN"); ok {
		// 		log.Debug("Reading Vault token from the environment (VAULT_TOKEN)")
		// 		viper.BindEnv("vault-token", "VAULT_TOKEN")
		// 	} else if _, err := os.Stat(vaultHandler.Config.TokenFile); err == nil {
		// 		log.Debug("Reading Vault token from (", vaultHandler.Config.TokenFile, ")")
		// 		token, err := ioutil.ReadFile(vaultHandler.Config.TokenFile)
		// 		check(err)
		// 		vaultHandler.Config.Token = string(token)
		// 	}
		//
		// 	// Set vaultHandler's config to have collected variables
		// 	err := viper.Unmarshal(&vaultHandler.Config)
		// 	if err != nil {
		// 		log.Fatal("unable to decode into struct: ", err)
		// 	}
		// },
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "login",
		Long:  "Stuff",
		Run: func(cmd *cobra.Command, args []string) {
			v.Login()
		},
	}

	v.api.BindCommand(vaultCmd, rootCmd)
	v.api.BindCommand(loginCmd, vaultCmd)
}

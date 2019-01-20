package cmd

import (
	"github.com/readytalk/stim/pkg/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var vaultHandler *vault.Client

func init() {
	rootCmd.AddCommand(vaultCmd)
	vaultCmd.AddCommand(vaultLoginCmd)

	// Get the Vault Address from the environment and command line
	vaultCmd.PersistentFlags().StringP("address", "a", "http://127.0.0.1:8200", "Vault URL")
	viper.BindEnv("vault-address", "VAULT_ADDR")
	viper.BindPFlag("vault-address", vaultCmd.PersistentFlags().Lookup("address"))
}

var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Vault helper",
	Long:  "Vault LDAP login and AWS access",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		vault.SetLogger(log) // Setup logging

		vaultHandler = vault.NewVault() // Setup a vault object

		// Set vaultHandler's config to have collected variables
		err := viper.Unmarshal(&vaultHandler.Config)
		if err != nil {
			log.Fatal("unable to decode into struct: ", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var vaultLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "login",
	Long:  "Stuff",
	Run: func(cmd *cobra.Command, args []string) {
    vaultHandler.Setup()
		// vaultHandler.Login()
	},
}

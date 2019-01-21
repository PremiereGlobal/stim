package vault

import (
	"github.com/readytalk/stim/stim"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (v *Vault) BindStim(s *stim.Stim) {
	v.stim = s
}

// The Command function is required for every stimpak
// This function sets up the cli command parameters and returns the command
func (v *Vault) Command(viper *viper.Viper) *cobra.Command {

	// var vaultCmd = &cobra.Command{
	// 	Use:   "vault",
	// 	Short: "Vault helper",
	// 	Long:  "Vault LDAP login and AWS access",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 	},
	// }

	var vaultCmd = &cobra.Command{
		Use:   "vault",
		Short: "Vault helper",
		Long:  "Vault LDAP login and AWS access",
		// PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// vault.SetLogger(log) // Setup logging
		//
		// vaultHandler = vault.NewVault() // Setup a vault object
		//
		// // Set vaultHandler's config to have collected variables
		// err := viper.Unmarshal(&vaultHandler.Config)
		// if err != nil {
		// 	log.Fatal("unable to decode into struct: ", err)
		// }
		// },
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// Get the Vault Address from the environment and command line
	vaultCmd.PersistentFlags().StringP("address", "a", "", "Vault URL")
	viper.BindEnv("vault-address", "VAULT_ADDR")
	viper.BindPFlag("vault-address", vaultCmd.PersistentFlags().Lookup("address"))

	// var loginCmd = &cobra.Command{
	// 	Use:   "login",
	// 	Short: "login",
	// 	Long:  "Stuff",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		v.Login()
	// 	},
	// }

	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "login",
		Long:  "Stuff",
		Run: func(cmd *cobra.Command, args []string) {
			v.Login()
			// vaultHandler.Setup()
			// vaultHandler.Login()
		},
	}

	v.stim.BindCommand(loginCmd, vaultCmd)

	return vaultCmd
}

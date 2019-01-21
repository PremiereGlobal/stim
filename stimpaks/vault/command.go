package vault

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// The Command function is required for every stimpak
// This function sets up the cli command parameters and returns the command
func (v *Vault) Command(viper *viper.Viper) *cobra.Command {

	var vaultCmd = &cobra.Command{
		Use:   "vault",
		Short: "Vault helper",
		Long:  "Vault LDAP login and AWS access",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	// Get the Vault Address from the environment and command line
	vaultCmd.PersistentFlags().StringP("address", "a", "http://127.0.0.1:8200", "Vault URL")
	viper.BindEnv("vault-address", "VAULT_ADDR")
	viper.BindPFlag("vault-address", vaultCmd.PersistentFlags().Lookup("address"))

	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "login",
		Long:  "Stuff",
		Run: func(cmd *cobra.Command, args []string) {
			v.Login()
		},
	}

	v.stim.BindCommand(loginCmd, vaultCmd)

	return vaultCmd
}

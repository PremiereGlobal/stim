package vault

import (
	"github.com/readytalk/stim/stim"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (v *Vault) BindStim(s *stim.Stim) {
	v.stim = s
}

// Command is required for every stimpack
// This function sets up the cli command parameters and returns the command
func (v *Vault) Command(viper *viper.Viper) *cobra.Command {

	var vaultCmd = &cobra.Command{
		Use:   "vault",
		Short: "Vault helper",
		Long:  "Vault LDAP login and AWS access",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Get the Vault Address from the environment and command line
	vaultCmd.PersistentFlags().StringP("address", "a", "", "Vault URL")
	viper.BindEnv("vault-address", "VAULT_ADDR")
	viper.BindPFlag("vault-address", vaultCmd.PersistentFlags().Lookup("address"))

	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "login to Vault",
		Long:  "Login into Vault using LDAP",
		Run: func(cmd *cobra.Command, args []string) {
			v.Login()
		},
	}
	v.stim.BindCommand(vaultCmd, loginCmd)

	// var awsCmd = &cobra.Command{
	// 	Use:   "aws",
	// 	Short: "aws login",
	// 	Long:  "Create AWS credentials",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		v.AWS()
	// 	},
	// }
	// v.stim.BindCommand(vaultCmd, awsCmd)

	return vaultCmd
}

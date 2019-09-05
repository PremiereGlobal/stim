package deploy

import (
	"github.com/PremiereGlobal/stim/stim"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (d *Deploy) BindStim(s *stim.Stim) {
	d.stim = s
}

// Command is required for every stimpack
// This function sets up the cli command parameters and returns the command
func (d *Deploy) Command(viper *viper.Viper) *cobra.Command {
	var deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy helper",
		Long:  "Deployment helper using Vault + Kubernetes + Helm",
		Run: func(cmd *cobra.Command, args []string) {
			d.Run()
		},
	}

	deployCmd.PersistentFlags().StringP("deploy-file", "f", "", "Deployment file")
	viper.BindPFlag("deploy.file", deployCmd.PersistentFlags().Lookup("deploy-file"))
	deployCmd.PersistentFlags().StringP("environment", "e", "", "Environment to deploy to")
	viper.BindPFlag("deploy.environment", deployCmd.PersistentFlags().Lookup("environment"))
	deployCmd.PersistentFlags().StringP("cluster", "c", "", "Cluster to deploy to")
	viper.BindPFlag("deploy.cluster", deployCmd.PersistentFlags().Lookup("cluster"))

	// Get the Vault Address from the environment and command line
	// vaultCmd.PersistentFlags().StringP("address", "a", "", "Vault URL")
	// viper.BindEnv("vault-address", "VAULT_ADDR")
	// viper.BindPFlag("vault-address", vaultCmd.PersistentFlags().Lookup("address"))

	// var loginCmd = &cobra.Command{
	// 	Use:   "login",
	// 	Short: "login to Vault",
	// 	Long:  "Login and obtain a token from Vault",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		v.Login()
	// 	},
	// }
	//
	// loginCmd.Flags().StringP("token-duration", "i", "", "Set token expiration for given duration. Example '8h'")
	// viper.BindPFlag("vault-initial-token-duration", loginCmd.Flags().Lookup("token-duration"))

	// v.stim.BindCommand(loginCmd, vaultCmd)
	return deployCmd
}

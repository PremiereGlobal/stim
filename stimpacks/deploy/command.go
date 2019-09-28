package deploy

import (
	"github.com/PremiereGlobal/stim/stim"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// BindStim creates the stim object within this stimpack
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
	deployCmd.PersistentFlags().StringP("instance", "i", "", "Instance to deploy to")
	viper.BindPFlag("deploy.instance", deployCmd.PersistentFlags().Lookup("instance"))
	deployCmd.PersistentFlags().StringP("method", "m", "docker", "Method to use for deployment.  Valid values are 'docker' or 'shell'")
	viper.BindPFlag("deploy.method", deployCmd.PersistentFlags().Lookup("method"))

	return deployCmd
}

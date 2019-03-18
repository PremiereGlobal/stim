package aws

import (
	"github.com/readytalk/stim/stim"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (a *Aws) BindStim(stim *stim.Stim) {
	a.stim = stim
	a.log = stim.GetLogger()
}

func (a *Aws) Command(viper *viper.Viper) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "aws",
		Short: "Interact with AWS",
		Long:  `Get credentials, etc.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "aws login",
		Long:  "Create AWS credentials",
		Run: func(cmd *cobra.Command, args []string) {
			err := a.Login()
			if err != nil {
				a.stim.Fatal(err)
			}
		},
	}
	a.stim.BindCommand(loginCmd, cmd)

	loginCmd.Flags().BoolP("source", "s", false, "output env source for current shell")
	viper.BindPFlag("env-source", loginCmd.Flags().Lookup("source"))

	loginCmd.Flags().BoolP("web", "w", false, "STS console web login")
	viper.BindPFlag("aws-web", loginCmd.Flags().Lookup("web"))

	loginCmd.Flags().StringP("mount", "m", "", "AWS Vault mount")
	viper.BindPFlag("aws-mount", loginCmd.Flags().Lookup("mount"))

	loginCmd.Flags().StringP("role", "r", "", "AWS Vault role")
	viper.BindPFlag("aws-role", loginCmd.Flags().Lookup("role"))

	return cmd
}

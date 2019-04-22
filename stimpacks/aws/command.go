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

	loginCmd.Flags().BoolP("web", "w", false, "Generate AWS web login (Default: launch URL)")
	viper.BindPFlag("aws-web", loginCmd.Flags().Lookup("web"))

	loginCmd.Flags().BoolP("output", "o", false, "Output URLs to console (don't launch URL)")
	viper.BindPFlag("aws-output", loginCmd.Flags().Lookup("output"))

	loginCmd.Flags().StringP("account", "a", "", "AWS Account")
	viper.BindPFlag("aws-account", loginCmd.Flags().Lookup("account"))

	loginCmd.Flags().StringP("role", "r", "", "AWS Vault role")
	viper.BindPFlag("aws-role", loginCmd.Flags().Lookup("role"))

	return cmd
}

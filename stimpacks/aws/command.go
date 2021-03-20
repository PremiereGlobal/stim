package aws

import (
	"github.com/PremiereGlobal/stim/stim"
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

	loginCmd.Flags().BoolP("use-profiles", "p", false, "Use profiles for storing credentials")
	viper.BindPFlag("aws.use-profiles", loginCmd.Flags().Lookup("use-profiles"))

	loginCmd.Flags().BoolP("default-profile", "d", false, "If --use-profiles is set, also set as [default] profile")
	viper.BindPFlag("aws.default-profile", loginCmd.Flags().Lookup("default-profile"))

	loginCmd.Flags().StringP("ttl", "t", "8h", "Time-to-live for AWS credentials")
	viper.BindPFlag("aws.ttl", loginCmd.Flags().Lookup("ttl"))

	loginCmd.Flags().StringP("web-ttl", "b", "1h", "Time-to-live for AWS web console access (min 15m, max 36h)")
	viper.BindPFlag("aws.web-ttl", loginCmd.Flags().Lookup("web-ttl"))

	loginCmd.Flags().BoolP("filter-by-token", "", false, "Show accounts and roles according to Vault token capabilities")
	viper.BindPFlag("aws.login.filter-by-token", loginCmd.Flags().Lookup("filter-by-token"))

	return cmd
}

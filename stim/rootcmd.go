package stim

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func initRootCommand(viper *viper.Viper) *cobra.Command {
	viper.SetEnvPrefix("stim")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	var cmd = &cobra.Command{
		Use:     "stim",
		Short:   "stim (stimulation delivery package) - Use your psychoactive hyperstimulants responsibly.",
		Long:    "Speeding up development with glue that brings tools together.",
		Example: "  To get help on a command:\n  stim vault help\n  To use a different config:\n  stim --config /tmp/config.yaml vault login",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Set root-level flags
	cmd.PersistentFlags().String("config", defaultStimConfigFilePath, "config file (default is "+defaultStimConfigFilePath+")")
	viper.BindPFlag("config-file", cmd.PersistentFlags().Lookup("config"))
	cmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	viper.BindPFlag("verbose", cmd.PersistentFlags().Lookup("verbose"))
	cmd.PersistentFlags().BoolP("noprompt", "x", false, "Do not prompt for input. Will default to true for Jenkin builds.")
	viper.BindPFlag("noprompt", cmd.PersistentFlags().Lookup("noprompt"))
	cmd.PersistentFlags().StringP("auth-method", "", "", "Default authentication method (ex: ldap, github, etc.)")
	viper.BindPFlag("auth.method", cmd.PersistentFlags().Lookup("auth-method"))
	cmd.PersistentFlags().BoolP("is-automated", "", false, "Error on anything that needs to prompt and was not passed in as an ENV var or command flag")
	viper.BindPFlag("is-automated", cmd.PersistentFlags().Lookup("is-automated"))

	// Set some defaults
	viper.SetDefault("vault-timeout", 15)

	return cmd
}

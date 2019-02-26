package stim

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (stim *Stim) rootCommand(viper *viper.Viper) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "stim",
		Short: "stim (stimulation delivery package) - Use your psychoactive hyperstimulants responsibly.",
		Long:  "Speeding up development with glue that brings tools together.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Set root-level flags
	cmd.PersistentFlags().String("config", "", "config file (default is $HOME/.stim.yaml)")
	viper.BindPFlag("config-file", cmd.PersistentFlags().Lookup("config"))
	cmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	viper.BindPFlag("verbose", cmd.PersistentFlags().Lookup("verbose"))
	cmd.PersistentFlags().BoolP("noprompt", "x", false, "Do not prompt for input. Will default to true for Jenkin builds.")
	viper.BindPFlag("noprompt", cmd.PersistentFlags().Lookup("noprompt"))

	// Set some defaults
	viper.SetDefault("vault-timeout", 15)

	return cmd
}

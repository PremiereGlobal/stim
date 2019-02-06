package stim

import (
  "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (stim *Stim) rootCommand(viper *viper.Viper) *cobra.Command {

  homeDir , err := homedir.Dir()
	var cmd = &cobra.Command{
		Use:   "stim",
		Short: "stim (stimulation delivery package) - Use your psychoactive hyperstimulants responsibly.",
		Long:  "Speeding up development with glue that brings tools together.",
    Example: "  To get help on a command:\n  stim vault help\n  To use a different config:\n  stim --config /tmp/config.yaml vault login",
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
  if homeDir == "" {
    if err != nil {
      stim.log.Debug("Could not get the home dir:", err)
    }
  } else {
    //TODO: we need to fix this for windows
    stim.config.SetDefault("config-file", homeDir+"/.stim/config.yaml")
    stim.config.SetDefault("homedir", homeDir)
  }

	return cmd
}

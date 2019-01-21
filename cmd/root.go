package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "stim",
	Short: "stim (stimulation delivery package) - Use your psychoactive hyperstimulants responsibly.",
	Long:  "Speeding up development with glue that brings tools together.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Command(viper *viper.Viper) *cobra.Command {

	// Set root-level flags
	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.stim.yaml)")
	viper.BindPFlag("config-file", rootCmd.PersistentFlags().Lookup("config"))
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	rootCmd.PersistentFlags().BoolP("noprompt", "x", false, "Do not prompt for input. Will default to true for Jenkin builds.")
	viper.BindPFlag("noprompt", rootCmd.PersistentFlags().Lookup("noprompt"))

	return rootCmd

}

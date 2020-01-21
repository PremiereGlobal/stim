package stim

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

func (stim *Stim) initRootCommand() {

	var cmd = &cobra.Command{
		Use:     "stim",
		Short:   "stim (stimulation delivery package) - Use your psychoactive hyperstimulants responsibly.",
		Long:    "Speeding up development with glue that brings tools together.",
		Example: "  To get help on a command:\n  stim vault help\n  To use a different config:\n  stim --config /tmp/config.yaml vault login",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// The default path for stim resources is in the ${HOME}/.stim directory
	// Users can override this with the STIM_PATH env
	// If the home directory cannot be determined, use a temp directory as the default
	defaultStimPath := ""
	home, err := homedir.Dir()
	if err != nil {
		defaultStimPath = filepath.Join(os.TempDir(), ".stim")
	} else {
		defaultStimPath = filepath.Join(home, ".stim")
	}
	stim.config.SetDefault("path", defaultStimPath)

	// The default cache path for stim resources
	// Can be overridden by a line in the config file or the env var STIM_CACHE_PATH
	defaultStimCachePath := filepath.Join(defaultStimPath, "cache")
	stim.config.SetDefault("cache-path", defaultStimCachePath)

	// Set root-level flags
	defaultStimConfigFilePath := filepath.Join(defaultStimPath, "config.yaml")
	cmd.PersistentFlags().String("config", defaultStimConfigFilePath, "config file (default is "+defaultStimConfigFilePath+")")
	stim.config.BindPFlag("config-file", cmd.PersistentFlags().Lookup("config"))

	cmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	stim.config.BindPFlag("verbose", cmd.PersistentFlags().Lookup("verbose"))
	cmd.PersistentFlags().BoolP("noprompt", "x", false, "Do not prompt for input. Will default to true for Jenkin builds.")
	stim.config.BindPFlag("noprompt", cmd.PersistentFlags().Lookup("noprompt"))
	cmd.PersistentFlags().StringP("auth-method", "", "", "Default authentication method (ex: ldap, github, etc.)")
	stim.config.BindPFlag("auth.method", cmd.PersistentFlags().Lookup("auth-method"))
	cmd.PersistentFlags().BoolP("is-automated", "", false, "Error on anything that needs to prompt and was not passed in as an ENV var or command flag")
	stim.config.BindPFlag("is-automated", cmd.PersistentFlags().Lookup("is-automated"))

	// Set some defaults
	stim.config.SetDefault("vault-timeout", 15)

	stim.rootCmd = cmd
}

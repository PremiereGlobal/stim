package cmd

import (
	homedir "github.com/mitchellh/go-homedir"
	"github.com/readytalk/stim/api"
	"github.com/readytalk/stim/stimpaks/pagerduty"
	"github.com/readytalk/stim/stimpaks/vault"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "stim",
	Short: "stim (stimulation delivery package) - Use your psychoactive hyperstimulants responsibly.",
	Long:  "Speeding up development with glue that brings tools together.",
	Run: func(cmd *cobra.Command, args []string) {
		// Show help if no sub-commands are specified
		log.Debug("Running root command")
		cmd.Help()
	},
}

// Version is normally set at build time
var version string
var config *viper.Viper

var configFile string
var log *logrus.Logger

func check(e error) { // This helper will streamline our error checks below.
	if e != nil {
		log.Error(e)
	}
}

func init() {

	// Set version for local testing if not set by build system
	if version == "" {
		version = "local"
	}

	// Initialize logger
	log = logrus.New()

	// Initialize viper (config)
	config = viper.New()

	// Initialize API
	log.Info("Initializing API")
	api := api.New(config)
	api.BindLogger(log)

	// Set root-level flags
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.stim.yaml)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
	config.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	rootCmd.PersistentFlags().BoolP("noprompt", "x", false, "Do not prompt for input. Will default to true for Jenkin builds.")
	config.BindPFlag("noprompt", rootCmd.PersistentFlags().Lookup("noprompt"))

	// Sets the passed functions to be run when each command's Execute method is called.
	cobra.OnInitialize(initConfig)

	// Initialize Pagerduty stimpak
	pagerduty := pagerduty.New(api)
	pagerduty.Bind(rootCmd)

	// Initialize Vault stimpak
	vault := vault.New(api)
	vault.Bind(rootCmd)
}

func initConfig() {

	config.SetConfigType("yaml") // set the config file type

	// Don't forget to read config either from CfgFile or from home directory!
	if configFile != "" {
		config.SetConfigFile(configFile) // Use config file from the flag.
	} else {
		home, err := homedir.Dir() // Find home directory.
		check(err)

		config.AddConfigPath(home)
		config.SetConfigName(".stim")
	}

	err := config.ReadInConfig()
	check(err)

	configFile = config.ConfigFileUsed()
	log.Info("config loaded")
	// Set log level
	if config.GetBool("verbose") == true {
		log.SetLevel(logrus.DebugLevel)
		log.Debug("Stim version: ", version)
		log.Debug("Debug log level set")
	}

	log.Debug("Using config file: ", configFile)

	if config.Get("noprompt") == false && isAutomated() {
		log.Debug("Detected automation. Setting --noprompt")
		config.Set("noprompt", true)
	}
}

func isAutomated() bool {
	if os.Getenv("JENKINS_URL") == "" {
		return false
	}
	return true
}

func Execute() {
	err := rootCmd.Execute()
	check(err)
}

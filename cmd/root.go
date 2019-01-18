package cmd

import (
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
  "github.com/sirupsen/logrus"
  homedir "github.com/mitchellh/go-homedir"
  "github.com/bartlettc22/stim-pkgs/pkg/notify"
)

var rootCmd = &cobra.Command {
  Use:   "stim",
  Short: "stim (stimulation delivery package) - Use your psychoactive hyperstimulants responsibly.",
  Long: "Speeding up development with glue that brings tools together.",
  PersistentPreRun: func(cmd *cobra.Command, args []string) {
  },
  Run: func(cmd *cobra.Command, args []string) {
  },
}

// Version is normally set at build time
var version string

var configFile string
var log *logrus.Logger
var Verbose bool

func check(e error) { // This helper will streamline our error checks below.
  if e != nil {
      log.Error(e)
  }
}

func init() {

  // Set version for local testing if not set by build system
  if version == ""  {
    version = "local"
  }

  log = logrus.New()
  cobra.OnInitialize(initConfig)
  rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.stim.yaml)")
  rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
  viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

  n := notify.New()
  n.BindCommand(rootCmd)
  n.BindLogger(log)
}

func initConfig() {
  viper.SetConfigType("yaml") // set the config file type

  // Don't forget to read config either from CfgFile or from home directory!
  if configFile != "" {
    viper.SetConfigFile(configFile)  // Use config file from the flag.
  } else {
    home, err := homedir.Dir() // Find home directory.
    check(err)

    viper.AddConfigPath(home)
    viper.SetConfigName(".stim")
  }

  err := viper.ReadInConfig()
  check(err)

  configFile = viper.ConfigFileUsed()

  if viper.GetBool("verbose") == true {
    log.SetLevel(logrus.DebugLevel)
    log.Debug("Debug level set")
  }

  log.Debug("Using config file: ", configFile)
}

func Execute() {
  err := rootCmd.Execute()
  check(err)
}

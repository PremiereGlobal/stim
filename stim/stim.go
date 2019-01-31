package stim

import (
	"github.com/readytalk/stim/pkg/vault"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"os"
	"os/user"
)

var version string

type Stim struct {
	config    *viper.Viper
	rootCmd   *cobra.Command
	log       *logrus.Logger
	stimpacks []*Stimpack
	version   string
	vault     *vault.Vault
}

var stim *Stim

func New() *Stim {

	// Create
	stim := &Stim{}

	// Initialize logger
	stim.log = logrus.New()

	// Initialize viper (config)
	stim.config = viper.New()

	// Set version for local testing if not set by build system
	if version == "" {
		stim.version = "local"
	} else {
		stim.version = version
	}

	stim.rootCmd = stim.rootCommand(stim.config)

	return stim
}

func (stim *Stim) Execute() {
	cobra.OnInitialize(stim.commandInit)
	err := stim.rootCmd.Execute()
	stim.Fatal(err)
}

func (stim *Stim) commandInit() {

	// Load a config file (if present)
	loadConfigErr := stim.loadConfigFile()

	// Set log level, this is done as early as possible so we can start using it
	if stim.GetConfigBool("verbose") == true {
		stim.log.SetLevel(logrus.DebugLevel)
		stim.log.Debug("Stim version: ", stim.version)
		stim.log.Debug("Debug log level set")
	}

	if loadConfigErr == nil {
		stim.log.Debug("Using config file: ", stim.config.ConfigFileUsed())
	} else {
		stim.log.Debug("No config file loaded")
		stim.log.Debug(loadConfigErr)
	}
}

func (stim *Stim) Vault() *vault.Vault {

	if stim.vault == nil {

		stim.log.Debug("Stim-Vault: Creating")

		username := stim.GetConfig("vault-username")
		if username == "" {
			var err error
			username, err = stim.User()
			if err != nil {
				stim.log.Fatal("Stim-vault: ", err)
			}
		}

		vault, err := vault.New(&vault.Config{
			Address:  stim.GetConfig("vault-address"),
			Noprompt: stim.GetConfigBool("noprompt") == false && stim.IsAutomated(),
			Logger:   stim.log,
			Username: username,
		})
		if err != nil {
			stim.log.Fatal("Stim-Vault: Error Initializaing: ", err)
		}

		// Update the user set in local configs to make new logins friendly
		err = stim.UpdateVaultUser(vault.GetUser())
		if err != nil {
			stim.log.Fatal("Stim-Vault: Error Updating username in configuration file: ", err)
		}

		stim.vault = vault
	}

	return stim.vault
}

func (stim *Stim) BindCommand(command *cobra.Command, parentCommand *cobra.Command) {
	parentCommand.AddCommand(command)
}

func (stim *Stim) User() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	return user.Username, nil
}

func (stim *Stim) UpdateVaultUser(username string) error {
	if username != stim.GetConfig("vault-username") {
		stim.Set("vault-username", username)
		err := stim.UpdateConfigFileKey("vault-username", username)
		if err != nil {
			return err
		}
	}

	return nil
}

func (stim *Stim) IsAutomated() bool {
	if os.Getenv("JENKINS_URL") == "" {
		return false
	}
	return true
}

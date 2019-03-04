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
		stim.log.Warn("Issue loading config file use -verbose for more info")
		stim.log.Debug(loadConfigErr)
	}
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

// UpdateVaultUser updates the user's stim config file with given username
// This username will be the default option when authenticating against Vault
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

// IsAutomated simply guesses if a build is invoking this code
// Used to disable user input prompts
func (stim *Stim) IsAutomated() bool {
	if os.Getenv("JENKINS_URL") == "" {
		return false
	}
	return true
}

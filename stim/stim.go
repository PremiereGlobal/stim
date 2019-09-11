package stim

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/PremiereGlobal/stim/pkg/stimlog"
	"github.com/PremiereGlobal/stim/pkg/vault"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var defaultStimConfigFilePath string
var version string

func init() {
	home, err := homedir.Dir()
	if err != nil {
		defaultStimConfigFilePath = filepath.Join(os.TempDir(), ".stim", "config.yaml")
	} else {
		defaultStimConfigFilePath = filepath.Join(home, ".stim", "config.yaml")
	}
	// Set version for local testing if not set by build system
	lv := "local"
	if version == "" {
		version = lv
	}
}

type Stim struct {
	config    *viper.Viper
	rootCmd   *cobra.Command
	log       stimlog.StimLogger
	stimpacks []*Stimpack
	vault     *vault.Vault
}

//New gets the Stim struct, which is treated like a singleton so you will get the same one
//as everywhere when this is called
func New() *Stim {
	log := stimlog.GetLogger()
	log.ForceFlush(true)
	config := viper.New()
	root := initRootCommand(config)
	return &Stim{log: log, config: config, rootCmd: root}
}

//GetLogger for Stim
func (stim *Stim) GetLogger() stimlog.StimLogger {
	return stim.log
}

func (stim *Stim) Execute() {
	cobra.OnInitialize(stim.commandInit)
	err := stim.rootCmd.Execute()
	stim.Fatal(err)
}

func (stim *Stim) commandInit() {
	// Load a config file (if present)
	loadConfigErr := stim.configLoadConfigFile()
	stim.configInitDefaultValues()

	if !stim.ConfigGetBool("logging.file.disable") {
		lfp := stim.ConfigGetString("logging.file.path")
		if lfp == "" {
			sh, err := stim.ConfigGetStimConfigDir()
			if err != nil {
				stim.log.Warn("Could not find stim config dir path, not creating log file")
			} else {
				lfp = filepath.Join(sh, "stim.log")
			}
		}
		if lfp != "" {
			stim.log.AddLogFile(lfp, stimlog.DebugLevel)
		}
	}
	// Set log level, this is done as early as possible so we can start using it
	if stim.ConfigGetBool("verbose") == true {
		// stim.log.SetLevel(logrus.DebugLevel)
		stim.log.SetLevel(stimlog.DebugLevel)
		stim.log.Debug("Stim version: {}", version)
		stim.log.Debug("Debug log level set")
	}
	if stim.IsAutomated() {
		stim.log.Info("Running in automated way")
	}

	if loadConfigErr == nil {
		stim.log.Debug("Using config file: {}", stim.config.ConfigFileUsed())
	} else if !stim.IsAutomated() {
		stim.log.Warn("Issue loading config file use --verbose for more info")
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
	if username != stim.ConfigGetString("vault-username") {
		stim.ConfigSetString("vault-username", username)
		err := stim.ConfigSetRaw("vault-username", username)
		if err != nil {
			return err
		}
	}

	return nil
}

// Used to disable user input prompts
func (stim *Stim) IsAutomated() bool {
	if stim.ConfigGetBool("is-automated") || os.Getenv("JENKINS_URL") != "" {
		return true
	}
	return false
}

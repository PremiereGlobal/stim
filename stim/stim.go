package stim

import (
	"os"
	"os/user"
	"sync"

	"github.com/readytalk/stim/pkg/stimlog"
	"github.com/readytalk/stim/pkg/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Stim struct {
	config    *viper.Viper
	rootCmd   *cobra.Command
	log       stimlog.StimLogger
	stimpacks []*Stimpack
	version   string
	vault     *vault.Vault
}

var version string
var stim *Stim

//New gets the Stim struct, which is treated like a singleton so you will get the same one
//as everywhere when this is called
func New() *Stim {
	if stim == nil {
		mu := sync.Mutex{}
		mu.Lock()
		if stim == nil {
			// Set version for local testing if not set by build system
			lv := "local"
			if version != "" {
				lv = version
			}
			log := stimlog.GetLogger()
      log.ForceFlush(true)
			config := viper.New()
			root := initRootCommand(config)
			stim = &Stim{log: log, config: config, rootCmd: root, version: lv}
		}
		mu.Unlock()
	}
	return stim
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
	loadConfigErr := stim.loadConfigFile()
	if !stim.GetConfigBool("disableLogFile") {
		lfp := stim.GetConfig("logFilePath")
		if lfp == "" {
			sh, err := stim.GetStimPath()
			if err != nil {
				stim.log.Warn("Could not find")
			} else {
				lfp = sh + "stim.log"
			}
		}
		if lfp != "" {
			stim.log.AddLogFile(lfp, stimlog.DebugLevel)
		}
	}
	// Set log level, this is done as early as possible so we can start using it
	if stim.GetConfigBool("verbose") == true {
		// stim.log.SetLevel(logrus.DebugLevel)
		stim.log.SetLevel(stimlog.DebugLevel)
		stim.log.Debug("Stim version: {}", stim.version)
		stim.log.Debug("Debug log level set")
	}

	if loadConfigErr == nil {
		stim.log.Debug("Using config file: {}", stim.config.ConfigFileUsed())
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

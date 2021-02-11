package stim

import (
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/PremiereGlobal/stim/pkg/stimlog"
	"github.com/PremiereGlobal/stim/pkg/vault"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/mod/semver"
)

var version string

func init() {

	// Set version for local testing if not set by build system
	lv := "v0.0.0-local"
	if version == "" {
		version = lv
	}
	if !semver.IsValid(version) {
		stimlog.GetLogger().Fatal("Bad Version:{}", version)
	}
}

type Stim struct {
	config    *viper.Viper
	rootCmd   *cobra.Command
	log       stimlog.StimLogger
	logConfig stimlog.StimLoggerConfig
	stimpacks []*Stimpack
	vault     *vault.Vault
}

//New gets the Stim struct, which is treated like a singleton so you will get the same one
//as everywhere when this is called
func New() *Stim {
	stim := &Stim{}
	stim.log = stimlog.GetLogger()
	stim.logConfig = stimlog.GetLoggerConfig()
	stim.logConfig.ForceFlush(true)
	stim.config = viper.New()
	stim.ConfigSetDefaultValues()
	stim.config.SetEnvPrefix("stim")
	stim.config.AutomaticEnv()
	stim.config.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	stim.initRootCommand()
	return stim
}

//GetLogger for Stim
func (stim *Stim) GetLogger() stimlog.StimLogger {
	return stim.log
}

func (stim *Stim) Execute() {
	defer stimlog.GetLoggerConfig().Flush()
	cobra.OnInitialize(stim.commandInit)
	err := stim.rootCmd.Execute()
	stim.Fatal(err)
}

func (stim *Stim) commandInit() {

	// Here we need to process certain config variables as this is the first time we have
	// access to them, including command line flags.
	// Particularly needed around flags that could change global pathing
	basePath := stim.config.GetString("path")
	if basePath == "" {
		home, err := homedir.Dir()
		if err != nil {
			basePath = filepath.Join(os.TempDir(), ".stim")
		} else {
			basePath = filepath.Join(home, ".stim")
		}
	}

	if stim.config.GetString("config-file") == "" {
		stim.config.Set("config-file", filepath.Join(basePath, "config.yaml"))
	}

	// Load a config file (if present)
	stim.configLoadConfigFile()

	// Now that we've loaded the config file, do one final check (in case path was set in the file)
	// If not set, use the basePath
	if stim.config.GetString("path") == "" {
		stim.config.Set("path", basePath)
	}

	if stim.config.GetString("cache-path") == "" {
		stim.config.Set("cache-path", filepath.Join(stim.config.GetString("path"), "cache"))
	}

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
			stim.logConfig.AddLogFile(lfp, stimlog.DebugLevel)
		}
	}
	// Set log level, this is done as early as possible so we can start using it
	if stim.ConfigGetBool("verbose") == true {
		stim.logConfig.SetLevel(stimlog.DebugLevel)
		stim.log.Debug("Stim version: {}", version)
		stim.log.Debug("Debug log level set")
	} else {
		// Set the default log level
		stim.logConfig.SetLevel(stimlog.InfoLevel)
	}
	if stim.IsAutomated() {
		stim.log.Info("Running in automated way")
	}

	stim.log.Debug("STIM_CONFIG_FILE: {}", stim.config.Get("config-file"))
	stim.log.Debug("STIM_PATH: {}", stim.config.Get("path"))
	stim.log.Debug("STIM_CACHE_PATH: {}", stim.config.Get("cache-path"))
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
	if username != stim.ConfigGetString("vault.username") {
		stim.ConfigSetString("vault.username", username)
		err := stim.ConfigSetRaw("vault.username", username)
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

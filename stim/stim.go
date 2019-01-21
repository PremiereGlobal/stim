package stim

import (
	"github.com/readytalk/stim/cmd"
	"github.com/readytalk/stim/pkg/pagerduty"
	"github.com/readytalk/stim/pkg/vault"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"os"
)

type Stim struct {
	config   *viper.Viper
	rootCmd  *cobra.Command
	log      *logrus.Logger
	stimpaks []*Stimpak
	version  string
}

// This is the interface for stimpaks
type Stimpak interface {
	Command(*viper.Viper) *cobra.Command
	Name() string
	BindStim(*Stim)
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
	if stim.version == "" {
		stim.version = "local"
	}

	stim.rootCmd = cmd.Command(stim.config)

	stim.commandInit()

	return stim
}

func (stim *Stim) AddStimpak(s Stimpak) {

	stim.log.Debug("Loading stimpak `", s.Name(), "`")
	s.BindStim(stim)
	// s.Stim = stim
	cmd := s.Command(stim.config)
	stim.rootCmd.AddCommand(cmd)
}

func (stim *Stim) Execute() {
	// rootCmd = cmd.rootCmd
	// cobra.OnInitialize(stim.commandInit)

	err := stim.rootCmd.Execute()
	stim.Fatal(err)
}

func (stim *Stim) commandInit() {

	// Load a config file (if present)
	loadConfigErr := stim.LoadConfigFile()

	// Set log level, this is done as early as possible so we can start using it
	if stim.config.GetBool("verbose") == true {
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

	if stim.config.Get("noprompt") == false && stim.isAutomated() {
		stim.log.Debug("Detected automation. Setting --noprompt")
		stim.config.Set("noprompt", true)
	}

}

func (stim *Stim) Pagerduty() *pagerduty.Pagerduty {
	stim.log.Debug("Stim-Pagerduty: Creating")
	vaultPath := stim.GetConfig("pagerduty.vault-apikey-path")
	vaultKey := stim.GetConfig("pagerduty.vault-apikey-key")
	stim.log.Debug("Stim-Pagerduty: Fetching Pagerduty API key from Vault `", vaultPath, "`")
	vault := stim.Vault()
	apikey, err := vault.GetSecretKey(vaultPath, vaultKey)
	if err != nil {
		stim.log.Fatal("Stim-Pagerduty: error getting API key from Vault: ", err)
	}
	pagerduty := pagerduty.New(apikey)
	return pagerduty
}

func (stim *Stim) Vault() *vault.Vault {

	stim.log.Debug("Stim-Vault: Creating")

	address := stim.GetConfig("vault-address")
	stim.log.Debug("Stim-Vault: Using Address ", address)

	vault, err := vault.New(&vault.Config{
		Address:  address,
		Noprompt: false,
	})
	// err := vault.InitClient()
	if err != nil {
		stim.log.Fatal("Stim-Vault: Error Initializaing: ", err)
	}
	return vault
}

func (stim *Stim) BindCommand(command *cobra.Command, parentCommand *cobra.Command) {
	parentCommand.AddCommand(command)
}

func (stim *Stim) isAutomated() bool {
	if os.Getenv("JENKINS_URL") == "" {
		return false
	}
	return true
}

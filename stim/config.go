package stim

import (
	"github.com/mitchellh/go-homedir"
)

func (stim *Stim) GetConfig(configKey string) string {
	configValue := stim.config.Get(configKey)
	if configValue != nil {
		return configValue.(string)
	}

	return ""
}

func (stim *Stim) GetConfigBool(configKey string) bool {
	configValue := stim.config.Get(configKey)
	if configValue != nil {
		return configValue.(bool)
	}

	return false
}

func (stim *Stim) Set(key string, value string) {
	stim.config.Set(key, value)
}

func (stim *Stim) UpdateConfigFile() error {
	err := stim.config.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}

func (stim *Stim) LoadConfigFile() error {

	// Set the config file type
	stim.config.SetConfigType("yaml")

	// Don't forget to read config either from CfgFile or from home directory!
	if configFile := stim.GetConfig("config-file"); configFile != "" {
		stim.config.SetConfigFile(configFile)
	} else {
		// Find home directory
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		stim.config.AddConfigPath(home)
		stim.config.SetConfigName(".stim")
	}

	err := stim.config.ReadInConfig()
	return err
}

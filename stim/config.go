package stim

import (
	"github.com/mitchellh/go-homedir"
	yaml "gopkg.in/yaml.v2"

	"io/ioutil"
	"os"
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

func (stim *Stim) UpdateConfigFileKey(key string, value string) error {
	config := make(map[string]interface{})

	f, err := ioutil.ReadFile(stim.config.ConfigFileUsed())
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(f, config)
	if err != nil {
		return err
	}

	config[key] = value

	f, err = yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(stim.config.ConfigFileUsed(), f, os.FileMode(0644))
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

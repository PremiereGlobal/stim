package stim

import (
	"github.com/mitchellh/go-homedir"

	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

func (stim *Stim) Get(configKey string) interface{} {
	configValue := stim.config.Get(configKey)
	if configValue != nil {
		return configValue
	}

	return nil
}

func (stim *Stim) GetConfig(configKey string) string {
	configValue := stim.config.GetString(configKey)
	return configValue
}

// GetConfigBool takes a config key and returns the boolean result
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

	var err error
	stimConfigFile := stim.config.ConfigFileUsed()
	if stimConfigFile == "" { // Will happen if the config doesn't exist
		stimConfigFile, err = stim.CreateConfigFile()
		if err != nil {
			return err
		}
	}

	// var f []byte
	f, err := ioutil.ReadFile(stimConfigFile)
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

	err = ioutil.WriteFile(stimConfigFile, f, os.FileMode(0644))
	if err != nil {
		return err
	}

	return nil
}

// CreateConfigFile will create the stim config file if it doesn't exist
// Used the frist time this code is ran so sub functions do not get errors when
// writting to the config.
func (stim *Stim) CreateConfigFile() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	stimConfigFile := home + "/.stim/config.yaml"

	dir, _ := path.Split(stimConfigFile)
	err = stim.CreateDirIfNotExist(dir)
	if err != nil {
		return "", err
	}
	newFile, err := os.Create(stimConfigFile)
	if err != nil {
		return "", err
	}
	newFile.Close()

	return stimConfigFile, nil
}

func (stim *Stim) CreateDirIfNotExist(dir string) error {
	stim.log.Debug("Creating config path: ", dir)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func (stim *Stim) loadConfigFile() error {

	// Set the config file type
	stim.config.SetConfigType("yaml")

	// Don't forget to read config either from CfgFile or from home directory!
	configFile := stim.GetConfig("config-file")
	_, err := os.Stat(configFile)
	if err != nil && !os.IsExist(err) {
		stim.log.Warn("No config file exits at :\"" + configFile + "\"")
		//If they passed in a custom path we might want to exit here
	}
	stim.config.SetConfigFile(configFile)
	confErr := stim.config.ReadInConfig()
	return confErr
}

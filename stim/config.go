package stim

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/PremiereGlobal/stim/pkg/utils"
	"github.com/imdario/mergo"
	yaml "gopkg.in/yaml.v3"
)

func (stim *Stim) ConfigGetRaw(configKey string) interface{} {
	var envCV interface{} = nil
	configValue := stim.config.Get(configKey)
	if strings.Contains(configKey, ".") {
		envCK := strings.ReplaceAll(configKey, ".", "-")
		envCV = stim.config.Get(envCK)
	}
	if envCV != nil {
		return envCV
	}
	if configValue != nil {
		return configValue
	}

	return nil
}

func (stim *Stim) ConfigGetString(configKey string) string {
	var envCV string = ""
	configValue := stim.config.GetString(configKey)
	if strings.Contains(configKey, ".") {
		envCK := strings.ReplaceAll(configKey, ".", "-")
		envCV = stim.config.GetString(envCK)
	}
	if envCV != "" {
		return envCV
	}
	return configValue
}

// GetConfigBool takes a config key and returns the boolean result
func (stim *Stim) ConfigGetBool(configKey string) bool {
	configValue := stim.ConfigGetRaw(configKey)
	if configValue != nil {
		switch v := configValue.(type) {
		case bool:
			return v
		case string:
			lc := strings.ToLower(v)
			return lc == "true" || lc == "t"
		default:
			stim.log.Warn("Unknown type for bool, type:{}, value of:{}", fmt.Sprintf("%T", v), v)
		}

	}
	return false
}

func (stim *Stim) ConfigHasValue(configKey string) bool {
	configValue := stim.config.Get(configKey)
	if strings.Contains(configKey, ".") {
		envCK := strings.ReplaceAll(configKey, ".", "-")
		envCV := stim.config.Get(envCK)
		if envCV != nil {
			return true
		}
	}
	if configValue != nil {
		return true
	}
	return false
}

func (stim *Stim) ConfigSetString(key string, value string) error {
	return stim.ConfigSetRaw(key, value)
}

func (stim *Stim) ConfigSetBool(key string, value bool) error {
	return stim.ConfigSetRaw(key, value)
}

func (stim *Stim) ConfigRemoveKey(key string) error {
	keys := []string{key}
	if strings.Contains(key, ".") {
		keys = strings.Split(key, ".")
	}
	config, err := stim.getConfigData()
	if err != nil {
		return err
	}
	cm, ok := config[keys[0]].(map[interface{}]interface{})
	if !ok {
		// Key doesn't exist, not thing to remove
		return nil
	}

	//This is super gross, not sure of a better way to deal with this
Main:
	for i, v := range keys[1:] {
		for cmk, cmv := range cm {
			if sv, ok := cmk.(string); ok {
				if i == len(keys)-2 {
					delete(cm, v)
					break Main
				} else {
					if sv == v {
						cm = cmv.(map[interface{}]interface{})
						break
					} else {
						return nil //We return nil because the entrie does not exit currently anyway
					}
				}
			}
		}
	}
	return stim.writeConfigData(config)
}

func (stim *Stim) ConfigSetRaw(key string, value interface{}) error {
	config := make(map[string]interface{})
	sc := make(map[string]interface{})
	if strings.Contains(key, ".") {
		keys := strings.Split(key, ".")
		bm := make(map[string]interface{})
		cm := bm
		for i, v := range keys {
			if i < len(keys)-1 {
				lm := make(map[string]interface{})
				cm[v] = lm
				cm = lm
			} else {
				cm[v] = value
			}
		}
		sc = bm
	} else {
		sc[key] = value
	}
	config, err := stim.getConfigData()
	if err != nil {
		return err
	}

	err = mergo.Merge(&config, sc)
	if err != nil {
		stim.log.Debug("Problem merging config:{}", err)
		return err
	}

	return stim.writeConfigData(config)
}

func (stim *Stim) writeConfigData(config map[string]interface{}) error {
	var err error
	stimConfigFile := stim.config.ConfigFileUsed()
	if stimConfigFile == "" { // Will happen if the config doesn't exist
		stimConfigFile, err = stim.ConfigGetStimConfigFile()
		if err != nil {
			return err
		}
	}
	f, err := yaml.Marshal(config)
	if err != nil {
		stim.log.Debug("Problem writing config yaml:{}", err)
		return err
	}
	err = ioutil.WriteFile(stimConfigFile, f, os.FileMode(0600))
	if err != nil {
		stim.log.Debug("Problem writing configfile:{}", err)
		return err
	}
	stim.getConfigData()
	return nil
}

func (stim *Stim) getConfigData() (map[string]interface{}, error) {
	config := make(map[string]interface{})
	var err error
	stimConfigFile := stim.config.ConfigFileUsed()
	if stimConfigFile == "" { // Will happen if the config doesn't exist
		stimConfigFile, err = stim.ConfigGetStimConfigFile()
		if err != nil {
			return nil, err
		}
	}

	f, err := ioutil.ReadFile(stimConfigFile)
	if err != nil {
		stim.log.Debug("Problem reading configfile:{}", err)
		return nil, err
	}
	err = yaml.Unmarshal(f, config)
	if err != nil {
		stim.log.Debug("Problem reading config yaml:{}", err)
		return nil, err
	}
	return config, err
}

func (stim *Stim) ConfigGetStimConfigDir() (string, error) {
	cfp, err := stim.ConfigGetStimConfigFile()
	if err != nil {
		return "", err
	}
	dir, _ := path.Split(cfp)
	return dir, nil
}

// ConfigGetStimConfigFile gets the current stim config file (creating it if necessary)
func (stim *Stim) ConfigGetStimConfigFile() (string, error) {

	cfp, err := filepath.Abs(stim.ConfigGetString("config-file"))

	// Create the config file if it doesn't exist
	// Probably want to abstract this to create default values as well
	err = utils.CreateFileIfNotExist(cfp, utils.UserOnlyMode)
	if err != nil {
		return "", err
	}

	return cfp, nil
}

func (stim *Stim) configLoadConfigFile() {

	stim.config.SetConfigType("yaml")
	configFile, err := stim.ConfigGetStimConfigFile()
	if err != nil {
		stim.log.Fatal("Problem accessing config file: {}", err)
	}

	stim.config.SetConfigFile(configFile)
	err = stim.config.ReadInConfig()
	if err != nil {
		stim.log.Fatal("Problem loading config file: {}", err)
	}

	// If the config file has a config-file entry remove it to avoid any sort
	// of circular reference.  This doesn't currently work so is commented out
	// to be dealt with in the future.  It doensn't work but ConfigRemoveKey only
	// removes it from the config file and not from the current stim.config
	// stim.ConfigRemoveKey("config-file") // old way
	// stim.ConfigRemoveKey("config.file") // new way
}

// ConfigGetStimCacheDir returns the stim cache directory
// subdir paramter optionally provides a subdirectory within the cache
func (stim *Stim) ConfigGetCacheDir(subDir string) string {

	cachePath := stim.ConfigGetString("cache-path")
	cacheSubPath := filepath.Join(cachePath, subDir)

	err := utils.CreateDirIfNotExist(cacheSubPath, utils.UserGroupMode)
	if err != nil {
		stim.log.Fatal("Error creating cache directory at {}", cacheSubPath)
	}

	return cacheSubPath
}

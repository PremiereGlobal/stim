package aws

import (
	"path/filepath"

	"github.com/PremiereGlobal/stim/pkg/utils"
	"github.com/go-ini/ini"
	"github.com/mitchellh/go-homedir"
)

// Profile is the base struct for creating new profiles
type Profile struct {
	AccessKeyID     string `ini:"aws_access_key_id"`
	SecretAccessKey string `ini:"aws_secret_access_key"`
}

// MapProfile gets an AWS profile from the user's credentials file and maps it,
// if found, to the given interface(s)
func (a *Aws) MapProfile(name string, profiles ...interface{}) error {

	// Get the profile configuration file path
	credentialPath, err := a.GetCredentialPath()
	if err != nil {
		return err
	}

	a.log.Debug("Mapping profile {} in {}", name, credentialPath)
	profileConfig, err := ini.Load(credentialPath)
	if err != nil {
		return err
	}

	// GetSection returns error if the section doesn't exist
	// If this is the case we'll just return here with no error
	section, err := profileConfig.GetSection(name)
	if err != nil {
		return nil
	}

	// Perform the map(s)
	for _, profile := range profiles {
		err = section.MapTo(profile)
		if err != nil {
			return err
		}
	}

	return nil
}

// SaveProfile saves the given profile(s)/section(s) to file, optionally also saving
// it as the default profile. The Profile type is requried and additional profile
// fields can be added with the additionalProfiles parameter
func (a *Aws) SaveProfile(name string, profile *Profile, setAsDefault bool, additionalProfiles ...interface{}) error {

	// Get the profile configuration file
	credentialPath, err := a.GetCredentialPath()
	if err != nil {
		return err
	}

	// Get the profile configuration file
	profileConfig, err := ini.Load(credentialPath)
	if err != nil {
		return err
	}

	// Write the profile
	writeSection(profileConfig, name, profile, additionalProfiles)

	// Overwrite the default profile, if set
	if setAsDefault {
		writeSection(profileConfig, "default", profile, additionalProfiles)
	}

	a.log.Debug("Saving profile {} in {}", name, credentialPath)
	profileConfig.SaveTo(credentialPath)

	return nil
}

// writeSection writes the given profile(s)/section(s) to the profileConfig
func writeSection(profileConfig *ini.File, name string, profile *Profile, additionalProfiles []interface{}) error {

	// Delete the section (if it exists)
	profileConfig.DeleteSection(name)

	// Create a new section to hold the new profile data
	section, err := profileConfig.NewSection(name)
	if err != nil {
		return err
	}

	// Reflect the main Profile into the newly created section
	err = section.ReflectFrom(profile)
	if err != nil {
		return err
	}

	// Reflect any additional profiles into the section
	for _, additionalProfile := range additionalProfiles {
		err = section.ReflectFrom(additionalProfile)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetCredentialPath gets the filepath to the credential path in the user's
// home directory
func (a *Aws) GetCredentialPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	credentialPath := filepath.FromSlash(home + "/.aws/credentials")
	err = utils.CreateFileIfNotExist(credentialPath)
	if err != nil {
		return "", err
	}

	return credentialPath, nil
}

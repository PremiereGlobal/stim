package aws

import (
	"github.com/go-ini/ini"
	"github.com/mitchellh/go-homedir"
	"path/filepath"
)

// MapProfile gets an AWS profile from the user's credentials file and maps it,
// if found, to the given interface
func (a *Aws) MapProfile(name string, profile interface{}) error {

	// Get the profile configuration file path
	credentialPath, err := a.GetCredentialPath()
	if err != nil {
		return err
	}

	a.log.Debug("Mapping profile " + name + " in " + credentialPath)
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

	// Perform the map
	err = section.MapTo(profile)
	if err != nil {
		return err
	}

	return nil
}

// SaveProfile saves the given profile/section to file, optionally also saving
// it as the default profile
func (a *Aws) SaveProfile(name string, profile interface{}, setAsDefault bool) error {

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
	writeSection(profileConfig, name, profile)

	// Overwrite the default profile, if set
	if setAsDefault {
		writeSection(profileConfig, "default", profile)
	}

	a.log.Debug("Saving profile " + name + " in " + credentialPath)
	profileConfig.SaveTo(credentialPath)

	return nil
}

// writeSection writes the given profile/section to the profileConfig
func writeSection(profileConfig *ini.File, name string, profile interface{}) error {

	// Delete the section (if it exists)
	profileConfig.DeleteSection(name)

	// Create a new section to hold the new profile data
	section, err := profileConfig.NewSection(name)
	if err != nil {
		return err
	}

	// Reflect profile into the newly created section
	err = section.ReflectFrom(profile)
	if err != nil {
		return err
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

	credentialPath := filepath.Join(home, ".aws/credentials")

	return credentialPath, nil
}

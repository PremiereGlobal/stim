package stim

import (
	"github.com/readytalk/stim/pkg/vault"
)

// Vault is the interface for Hashicorp Vault wrapper methods
// The main input is the vault-address
// Will prompt the user for their LDAP username and password
// Will update the user's ~/.vault-token file with a new token
func (stim *Stim) Vault() *vault.Vault {
	if stim.vault == nil {

		stim.log.Debug("Stim-Vault: Creating")

		username := stim.GetConfig("vault-username")
		if username == "" {
			var err error
			username, err = stim.User()
			if err != nil {
				stim.log.Fatal("Stim-vault: ", err)
			}
		}

		// Create the Vault object and pass in the needed address
		vault, err := vault.New(&vault.Config{
			Address:  stim.GetConfig("vault-address"), // Default is 127.0.0.1
			Noprompt: stim.GetConfigBool("noprompt") == false && stim.IsAutomated(),
			Log:      stim.log, // Pass in the global logger object
			Username: username, // If set in the configs, pass in user
		})
		if err != nil {
			stim.log.Fatal("Stim-Vault: Error Initializaing: ", err)
		}
		stim.vault = vault

		// Update the username set in local configs to make logins more friendly
		err = stim.UpdateVaultUser(vault.GetUser())
		if err != nil {
			stim.log.Fatal("Stim-Vault: Error Updating username in configuration file: ", err)
		}

		// If user wants, extend the token timeout
		if stim.vault.IsNewLogin() {
			renewTime := stim.GetConfig("vault-renew")
			if renewTime != "" {
				stim.vault.RenewToken(renewTime)
			}
		}
	}

	return stim.vault
}

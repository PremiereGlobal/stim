package stim

import (
	"github.com/PremiereGlobal/stim/pkg/vault"

	"time"
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

		// Note with ParseDuration: If you value is 28800 you will need to add an "s" at the end
		timeInDuration, err := time.ParseDuration(stim.GetConfig("vault-initial-token-duration"))
		if err != nil {
			stim.log.Warn("Stim-vault: ", err)
			timeInDuration = time.Duration(0)
		}

		// Create the Vault object and pass in the needed address
		vault, err := vault.New(&vault.Config{
			Address:              stim.GetConfig("vault-address"), // Default is 127.0.0.1
			Noprompt:             stim.GetConfigBool("noprompt") == false && stim.IsAutomated(),
			AuthPath:             stim.GetConfig("auth.method"),
			Username:             username, // If set in the configs, pass in user
			InitialTokenDuration: timeInDuration,
			Log:                  stim.log,
		})
		if err != nil {
			stim.log.Fatal(err)
		}
		stim.vault = vault

		// Update the username set in local configs to make logins more friendly
		err = stim.UpdateVaultUser(vault.GetUser())
		if err != nil {
			stim.log.Fatal("Stim-Vault: Error Updating username in configuration file: ", err)
		}
	}

	return stim.vault
}

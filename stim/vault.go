package stim

import (
	"github.com/PremiereGlobal/stim/pkg/vault"

	"time"
)

// Vault is the interface for Hashicorp Vault wrapper methods
// The main input is the vault-address
// Will update the user's ~/.vault-token file with a new token
func (stim *Stim) Vault() *vault.Vault {
	if stim.vault == nil {

		stim.log.Debug("Stim-Vault: Creating")

		username := stim.ConfigGetString("vault.username")

		// Note with ParseDuration: If you value is 28800 you will need to add an "s" at the end
		var timeInDuration time.Duration
		var err error

		timeInDuration = stim.ConfigGetDuration("vault.default-ttl")
		if timeInDuration == time.Duration(0) {
			vtd := stim.ConfigGetString("vault-initial-token-duration") //TODO: depreciated config should be removed
			if vtd != "" {
				timeInDuration, err = time.ParseDuration(vtd)
				if err != nil {
					stim.log.Warn("Stim-vault: bad duration value:{} caused error:{}", vtd, err)
					timeInDuration = time.Duration(0)
				}
			}
		}

		va := stim.ConfigGetString("vault.address")
		stim.log.Debug("Vault Address: ({})", va)

		skipUserPrompt := false
		skipUserPrompt = stim.ConfigGetBool("vault.skip-username-prompt")
		if !skipUserPrompt {
			skipUserPrompt = stim.ConfigGetBool("vault-username-skip-prompt") //TODO: depreciated config should be removed
		}

		// Create the Vault object and pass in the needed address
		vault, err := vault.New(&vault.Config{
			Address:              va, // Default is 127.0.0.1
			Noprompt:             stim.ConfigGetBool("noprompt") == false && stim.IsAutomated(),
			AuthPath:             stim.ConfigGetString("auth.method"),
			Username:             username, // If set in the configs, pass in user
			UsernameSkipPrompt:   skipUserPrompt,
			InitialTokenDuration: timeInDuration,
			Log:                  stim.log,
		})
		if err != nil {
			stim.log.Fatal(err)
		}
		stim.vault = vault

		// Update the username set in local configs to make logins more friendly
		err = stim.UpdateVaultUser(vault.GetUser())
		if err != nil && !stim.IsAutomated() {
			stim.log.Warn("Stim-Vault: Could not Update username in config file: {}", err)
		}
	}

	return stim.vault
}

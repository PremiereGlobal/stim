package stim

import (
	"github.com/PremiereGlobal/stim/pkg/pagerduty"
)

// Pagerduty returns a Pagerduty instance that is already authenticated
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
	pagerduty := pagerduty.New(apikey, stim.log)
	return pagerduty
}

package aws

import (
	// "github.com/readytalk/stim/pkg/log"

	"errors"
)

// GetCredentials will get the aws mount and role from the user
func (a *Aws) GetCredentials() (string, string, error) {
	if a.vault == nil {
		a.vault = a.stim.Vault()
	}

	mounts, err := a.vault.GetMounts("aws")
	a.stim.Fatal(err)

	vaultAccount := a.stim.GetConfig("aws-mount")
	if vaultAccount == "" && a.stim.IsAutomated() {
		a.stim.Fatal(errors.New("Vault aws mount not specified"))
	} else if vaultAccount == "" {
		vaultAccount, err = a.stim.PromptSearchList("Choose AWS account", mounts)
		a.stim.Fatal(err)
	}

	vaultRole := a.stim.GetConfig("aws-role")
	if vaultRole == "" && a.stim.IsAutomated() {
		a.stim.Fatal(errors.New("Vault aws role not specified"))
	} else if vaultRole == "" {
		vaultRole, err = a.stim.PromptListVault(vaultAccount+"/roles", "Select Role", "")
		a.stim.Fatal(err)
	}

	return vaultAccount, vaultRole, nil
}

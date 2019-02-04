package aws

import (
	"errors"
	"github.com/davecgh/go-spew/spew"
)

func (a *Aws) GetCredentials() {
	vault := a.stim.Vault()

	mounts, err := vault.GetMounts("aws")
	a.stim.Fatal(err)

	// k.stim.GetConfig("kube-service-account")
	vaultAccount := a.stim.GetConfig("aws.vault-mount")
	if vaultAccount == "" && a.stim.IsAutomated() {
		a.stim.Fatal(errors.New("Vault aws mount not specified"))
	} else if vaultAccount == "" {
		vaultAccount, err = a.stim.PromptSearchList("Choose AWS account", mounts)
		a.stim.Fatal(err)
	}

	//k.stim.GetConfig("kube-service-account")
	vaultRole := a.stim.GetConfig("aws.vault-role")
	if vaultRole == "" && a.stim.IsAutomated() {
		a.stim.Fatal(errors.New("Vault aws role not specified"))
	} else if vaultRole == "" {
		vaultRole, err = a.stim.PromptListVault(vaultAccount+"roles", "Select Role", "")
		a.stim.Fatal(err)
	}

	spew.Dump(vaultRole)
}

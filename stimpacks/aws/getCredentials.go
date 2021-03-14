package aws

import (
	"errors"
	"sort"
	"strings"
)

// GetCredentials will get the aws mount and role from the user
func (a *Aws) GetCredentials() (string, string, error) {
	if a.vault == nil {
		a.vault = a.stim.Vault()
	}

	vaultAccount := a.stim.ConfigGetString("aws.account")
	if vaultAccount == "" && a.stim.IsAutomated() {
		a.stim.Fatal(errors.New("Vault aws mount not specified"))
	} else if vaultAccount == "" {
		filteredAccounts, err := a.filterAccounts()
		vaultAccount, err = a.stim.PromptSearchList("Choose AWS account", filteredAccounts)
		a.stim.Fatal(err)
	}

	vaultRole := a.stim.ConfigGetString("aws.role")
	if vaultRole == "" && a.stim.IsAutomated() {
		a.stim.Fatal(errors.New("Vault aws role not specified"))
	} else if vaultRole == "" {
		filteredRoles, err := a.filterRoles(vaultAccount)
		vaultRole, err = a.stim.PromptList("Select Role", filteredRoles, "")
		a.stim.Fatal(err)
	}

	return vaultAccount, vaultRole, nil
}

func (a *Aws) filterAccounts() ([]string, error) {
	mounts, err := a.vault.GetMounts("aws")
	a.stim.Fatal(err)
	if !a.stim.ConfigGetBool("aws.filter-prompts") {
		return mounts, nil
	}

	var rolesPaths []string
	for _, mount := range mounts {
		rolesPaths = append(rolesPaths, mount+"/roles")
	}

	filteredPaths, err := a.vault.Filter(rolesPaths, []string{"list"})
	a.stim.Fatal(err)

	var filteredAccounts []string
	for _, path := range filteredPaths {
		account := strings.TrimSuffix(path, "/roles")
		filteredAccounts = append(filteredAccounts, account)
	}

	sort.Strings(filteredAccounts)
	return filteredAccounts, nil
}

func (a *Aws) filterRoles(vaultAccount string) ([]string, error) {
	roles, err := a.vault.ListSecrets(vaultAccount + "/roles")
	a.stim.Fatal(err)
	if !a.stim.ConfigGetBool("aws.filter-prompts") {
		return roles, nil
	}

	var credsPaths []string
	for _, role := range roles {
		credsPaths = append(credsPaths, vaultAccount+"/creds/"+role)
	}

	filteredPaths, err := a.vault.Filter(credsPaths, []string{"read"})
	a.stim.Fatal(err)

	var filteredRoles []string
	for _, path := range filteredPaths {
		role := strings.TrimPrefix(path, vaultAccount+"/creds/")
		filteredRoles = append(filteredRoles, role)
	}

	sort.Strings(filteredRoles)
	return filteredRoles, nil
}

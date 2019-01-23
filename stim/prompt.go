package stim

import (
	"github.com/manifoldco/promptui"
	"strings"
)

// Prompt for yes/no type question
// label: prompt label
// override: if set to true, this function will return true
// default: what will be used if nothing is entered
func (stim *Stim) PromptBool(label string, override bool, defaultvalue bool) (bool, error) {

	if override {
		stim.Debug("PromptString: Using override value of `true`")
		return true, nil
	}

	defaultstring := "n"
	if defaultvalue {
		defaultstring = "y"
	}
	label = label + " (y/n) [" + defaultstring + "]"

	prompt := promptui.Prompt{
		Label: label,
	}

	result, err := prompt.Run()
	if err != nil {
		return false, err
	}

	if result == "" || strings.ToLower(strings.TrimSpace(result))[0:1] == "y" {
		return true, nil
	}

	return false, nil
}

// Prompt for string
// label: prompt label
// override: if set, will use this value instead of prompting
// default: what will be used if nothing is entered
func (stim *Stim) PromptString(label string, override string, defaultvalue string) (string, error) {

	if override != "" {
		stim.Debug("PromptString: Using override value of `" + override + "`")
		return override, nil
	}

	defaultstring := ""
	if defaultvalue != "" {
		defaultstring = "[" + defaultvalue + "] "
	}
	label = label + " " + defaultstring + ""

	prompt := promptui.Prompt{
		Label: label,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	if result == "" {
		return defaultvalue, nil
	}

	return result, nil
}

// Prompt List
// label: prompt label
// override: if set, will use this value instead of prompting
// default: what will be used if nothing is entered
func (stim *Stim) PromptList(label string, list []string, override string) (string, error) {

	if override != "" {
		stim.Debug("PromptList: Using override value of `" + override + "`")
		return override, nil
	}

	prompt := promptui.Select{
		Label: label,
		Items: list,
		Size:  10,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

// Prompt List
// label: prompt label
// override: if set, will use this value instead of prompting
// default: what will be used if nothing is entered
func (stim *Stim) PromptListVault(vaultPath string, label string, override string) (string, error) {

	if override != "" {
		stim.Debug("PromptListVault: Using override value of `" + override + "`")
		return override, nil
	}

	vault := stim.Vault()
	list, err := vault.ListSecrets(vaultPath)
	if err != nil {
		return "", err
	}

	result, err := stim.PromptList(label, list, "")
	if err != nil {
		return "", err
	}

	return result, nil
}

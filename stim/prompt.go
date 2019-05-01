package stim

import (
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/manifoldco/promptui"
)

// stderr implements an io.WriteCloser that skips the terminal bell character
// (ASCII code 7), and writes the rest to os.Stderr. It's used to replace
// readline.Stdout, that is the package used by promptui to display the prompts.
type stderr struct{}

// Write implements an io.WriterCloser over os.Stderr, but it skips the terminal
// bell character.
func (s *stderr) Write(b []byte) (int, error) {
	if len(b) == 1 && b[0] == readline.CharBell {
		return 0, nil
	}
	return os.Stderr.Write(b)
}

// Close implements an io.WriterCloser over os.Stderr.
func (s *stderr) Close() error {
	return os.Stderr.Close()
}

func init() {
	readline.Stdout = &stderr{}
}

// PromptBool asks the user a yes/no question
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

// PromptString prompts the user to enter a string
func (stim *Stim) PromptString(label string, defaultvalue string) (string, error) {

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

// PromptList prompts the user to select from the list of string provided
// If override string is not empty it will be returned without
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

// PromptListVault uses a path from vault and prompts to select the list
// of secrets within that list.  Returns the value selected.
// If override string is not empty it will be returned without
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

// PromptSearchList takes a label, list of selectable values and prompts the user
// to select the results.  If override string is not empty it will be returned without
// prompting
func (stim *Stim) PromptSearchList(label string, list []string) (string, error) {

	searcher := func(input string, index int) bool {
		name := strings.Replace(strings.ToLower(list[index]), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:             label,
		Items:             list,
		Size:              10,
		Searcher:          searcher,
		StartInSearchMode: true,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

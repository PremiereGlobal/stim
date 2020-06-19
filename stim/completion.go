package stim

import (
	"fmt"
	"io"
	"os"
)

func (stim *Stim) GetCompletion(shell string) error {
	switch shell {
	case `bash`:
		stim.rootCmd.GenBashCompletion(os.Stdout)
	case `zsh`:
		stim.rootCmd.GenZshCompletion(os.Stdout)
		io.WriteString(os.Stdout, "\ncompdef _stim stim\n")
	default:
		return fmt.Errorf("Unknown shell: %s", shell)
	}

	return nil
}

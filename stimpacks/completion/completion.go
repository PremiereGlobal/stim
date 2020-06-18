package completion

import (
	"fmt"
	"github.com/PremiereGlobal/stim/stim"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Completion struct {
	name string
	stim *stim.Stim
}

func New() *Completion {
	return &Completion{name: "completion"}
}

func (c *Completion) Name() string {
	return c.name
}

func (c *Completion) BindStim(s *stim.Stim) {
	c.stim = s
}

func (c *Completion) Command(viper *viper.Viper) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "completion SHELL",
		Short: "Output shell completion for the given shell (bash or zsh)",
		Long: `Output shell completion for the given shell (bash or zsh)
The following ought to suffice for loading the Bash completions:
	source <(stim completion bash)

Zsh is more complicated because there is more than one completion engine for
Zsh. Try putting the completion output into a script (in your $fpath) and 
loading it with compinit.
	stim completion zsh > /path/to/script

		`,
		ValidArgs: []string{"bash", "zsh"},
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.stim.GetCompletion(args[0]); err != nil {
				fmt.Println(err)
			}
		},
	}

	return cmd
}

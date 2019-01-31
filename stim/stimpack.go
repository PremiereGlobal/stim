package stim

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// This is the interface for stimpacks
type Stimpack interface {
	Command(*viper.Viper) *cobra.Command
	Name() string
	BindStim(*Stim)
}

func (stim *Stim) AddStimpack(s Stimpack) {

	stim.log.Debug("Loading stimpack `", s.Name(), "`")
	s.BindStim(stim)
	cmd := s.Command(stim.config)
	stim.rootCmd.AddCommand(cmd)
}

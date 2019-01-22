package stim

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// This is the interface for stimpaks
type Stimpak interface {
	Command(*viper.Viper) *cobra.Command
	Name() string
	BindStim(*Stim)
}

func (stim *Stim) AddStimpak(s Stimpak) {

	stim.log.Debug("Loading stimpak `", s.Name(), "`")
	s.BindStim(stim)
	cmd := s.Command(stim.config)
	stim.rootCmd.AddCommand(cmd)
}

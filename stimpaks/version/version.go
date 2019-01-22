package version

import (
	"fmt"
	"github.com/readytalk/stim/stim"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Version struct {
	name string
	stim *stim.Stim
}

func New() *Version {
	return &Version{name: "version"}
}

func (v *Version) Name() string {
	return v.name
}

func (v *Version) BindStim(s *stim.Stim) {
	v.stim = s
}

func (v *Version) Command(viper *viper.Viper) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "version",
		Short: "Print the client version",
		Long:  `Print the client version`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("stim/%v\n", v.stim.GetVersion())
		},
	}

	return cmd
}

package discover

import (
	"fmt"
	"github.com/readytalk/stim/stim"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (d *Discover) BindStim(s *stim.Stim) {
	d.stim = s
}

func (d *Discover) Command(viper *viper.Viper) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "discover",
		Short: "Use to discover services, endpoints, etc.",
		Long:  "Use to discover services, endpoints, etc.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	var clusterCmd = &cobra.Command{
		Use:   "clusters",
		Short: "Use to discover services, endpoints, etc.",
		Long:  "Use to discover services, endpoints, etc.",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := d.DiscoverClusters()
			if err != nil {
				d.stim.Fatal(err)
			} else {
				fmt.Println(result)
			}
		},
	}

	// cmd.Flags().BoolP("clusters", "", false, "Display list of known Kubernetes clusters")
	// viper.BindPFlag("discover-clusters", cmd.Flags().Lookup("clusters"))

	d.stim.BindCommand(clusterCmd, cmd)

	return cmd
}

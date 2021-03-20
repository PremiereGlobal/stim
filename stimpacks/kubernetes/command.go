package kubernetes

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (k *Kubernetes) Command(viper *viper.Viper) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "kube",
		Short: "Used to config and interact with Kubernetes",
		Long:  "Used to config and interact with Kubernetes",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	var configCmd = &cobra.Command{
		Use:   "config",
		Short: "Create/modify a Kubernetes context",
		Long:  "Create/modify a Kubernetes context",
		Run: func(cmd *cobra.Command, args []string) {
			err := k.configureContext()
			if err != nil {
				k.stim.Fatal(err)
			}
		},
	}

	configCmd.Flags().StringP("cluster", "c", "", "Required. Name of cluster to config")
	viper.BindPFlag("kube-config-cluster", configCmd.Flags().Lookup("cluster"))
	configCmd.Flags().StringP("service-account", "s", "", "Required. Name of service account to use")
	viper.BindPFlag("kube-service-account", configCmd.Flags().Lookup("service-account"))
	configCmd.Flags().StringP("context", "t", "", "Optional. Name of context to set. Default is cluster name")
	viper.BindPFlag("kube-context", configCmd.Flags().Lookup("context"))
	configCmd.Flags().BoolP("current-context", "r", false, "Optional. Set to current context")
	viper.BindPFlag("kube-current-context", configCmd.Flags().Lookup("current-context"))
	configCmd.Flags().StringP("namespace", "n", "", "Optional. Name of default namespace")
	viper.BindPFlag("kube-config-namespace", configCmd.Flags().Lookup("namespace"))
	configCmd.Flags().StringP("cf", "", "", "Optional. Cluster regex filter")
	viper.BindPFlag("kube.config.cluster-filter", configCmd.Flags().Lookup("cf"))
	configCmd.Flags().StringP("saf", "", "", "Optional. Service Account regex filter")
	viper.BindPFlag("kube.config.service-account-filter", configCmd.Flags().Lookup("saf"))
	configCmd.Flags().BoolP("filter-by-token", "", false, "Optional. Show service accounts by Vault token capabilities")
	viper.BindPFlag("kube.config.filter-by-token", configCmd.Flags().Lookup("filter-by-token"))

	k.stim.BindCommand(configCmd, cmd)

	return cmd
}

package kubernetes

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (k *Kubernetes) Command(viper *viper.Viper) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "kube",
		Short: "Used to configure and interact with Kubernetes",
		Long:  "Used to configure and interact with Kubernetes",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	var configureCmd = &cobra.Command{
		Use:   "configure",
		Short: "Create/modify a Kubernetes context",
		Long:  "Create/modify a Kubernetes context",
		Run: func(cmd *cobra.Command, args []string) {
			err := k.configureContext()
			if err != nil {
				k.stim.Fatal(err)
			}
		},
	}

	configureCmd.Flags().StringP("cluster", "c", "", "Required. Name of cluster to configure")
	viper.BindPFlag("kube-configure-cluster", configureCmd.Flags().Lookup("cluster"))
	configureCmd.Flags().StringP("service-account", "s", "", "Required. Name of service account to use")
	viper.BindPFlag("kube-service-account", configureCmd.Flags().Lookup("service-account"))
	configureCmd.Flags().StringP("context", "t", "", "Optional. Name of context to set. Default is cluster name")
	viper.BindPFlag("kube-context", configureCmd.Flags().Lookup("context"))
	configureCmd.Flags().BoolP("current-context", "r", false, "Optional. Set to current context")
	viper.BindPFlag("kube-current-context", configureCmd.Flags().Lookup("current-context"))
	configureCmd.Flags().StringP("namespace", "n", "", "Optional. Name of default namespace")
	viper.BindPFlag("kube-config-namespace", configureCmd.Flags().Lookup("namespace"))

	k.stim.BindCommand(configureCmd, cmd)

	return cmd
}

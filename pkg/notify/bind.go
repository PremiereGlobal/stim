package notify

import (
	VaultApi "github.com/hashicorp/vault/api"
	"github.com/readytalk/stim/pkg/pagerduty"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (n *Notify) BindLogger(log *logrus.Logger) {
	n.log = log
}

func (n *Notify) BindVault(vault *VaultApi.Client) {
	n.vault = vault
}

func (n *Notify) BindCommand(parentCommand *cobra.Command) {

	var cmd = &cobra.Command{
		Use:   "notify",
		Short: "Sends notifications",
		Long:  `Sends notifications to platforms like Pagerduty and Slack`,
		Run: func(cmd *cobra.Command, args []string) {
			n.log.Debug("Running notify command")
			viper.Unmarshal(n.Options)
			n.Notify()
		},
	}

	cmd.Flags().StringP("type", "t", "", "Type of notification. Valid values are [\"generic\",\"jenkins-job\"].  Default is \"generic\"")
	viper.BindPFlag("notify-type", cmd.Flags().Lookup("type"))
	viper.SetDefault("notify-type", "generic")

	cmd.Flags().String("pagerduty", "", "Pagerduty service to send notification to")
	viper.BindPFlag("notify-pagerduty", cmd.Flags().Lookup("pagerduty"))

	cmd.Flags().String("slack", "", "Slack channel to send notification to")
	viper.BindPFlag("notify-slack", cmd.Flags().Lookup("slack"))

	cmd.Flags().StringP("status", "s", "", "Status of notification.  Valid values are [\"success\",\"failure\"]")
	viper.BindPFlag("notify-status", cmd.Flags().Lookup("status"))

	cmd.Flags().StringP("message", "m", "", "Message to send in notification.  Required if --type is \"generic\"")
	viper.BindPFlag("notify-message", cmd.Flags().Lookup("message"))

	cmd.PersistentFlags().StringP("vault-pagerduty-apikey-path", "", "", "Path in Vault for the Pagerduty API key.  Required if using Pagerduty.")
	viper.BindPFlag("vault-pagerduty-apikey-path", cmd.PersistentFlags().Lookup("vault-pagerduty-apikey-path"))

	cmd.PersistentFlags().StringP("vault-pagerduty-apikey-key", "", "", "Key in Vault path containing the Pagerduty API key.  Required if using Pagerduty.")
	viper.BindPFlag("vault-pagerduty-apikey-key", cmd.PersistentFlags().Lookup("vault-pagerduty-apikey-key"))

	// Add subcommand(s)
	p := pagerduty.New()
	p.BindLogger(n.log)
	p.BindVault(n.vault)
	p.BindCommand(cmd)

	parentCommand.AddCommand(cmd)
}

package notify

import (
  "github.com/spf13/cobra"
  "github.com/spf13/viper"
  "github.com/sirupsen/logrus"
)

func (n *Notify) BindCommand(parentCommand *cobra.Command) {

  var cmd = &cobra.Command{
    Use:   "notify",
    Short: "Sends notifications",
    Long:  `Sends notifications to platforms like Pagerduty and Slack`,
    Run: func(cmd *cobra.Command, args []string) {
      n.log.Debug("Running notify command")
      viper.Unmarshal(n.Options)
      n.log.Debug(n)
      n.Notify()
    },
  }

  cmd.Flags().StringP("type", "t", "", "Type of notification. Valid values are [\"generic\",\"jenkins-job\"].  Default is \"generic\"")
  viper.BindPFlag("notify-type", cmd.Flags().Lookup("type"))

  cmd.Flags().String("pagerduty", "", "Pagerduty service to send notification to")
  viper.BindPFlag("notify-pagerduty", cmd.Flags().Lookup("pagerduty"))

  cmd.Flags().String("slack", "", "Slack channel to send notification to")
  viper.BindPFlag("notify-slack", cmd.Flags().Lookup("slack"))

  cmd.Flags().StringP("status", "s", "", "Status of notification.  Valid values are [\"info\"]. Default is \"info\"")
  viper.BindPFlag("notify-status", cmd.Flags().Lookup("status"))

  parentCommand.AddCommand(cmd)
}

func (n *Notify) BindLogger(log *logrus.Logger) {
  n.log = log
}

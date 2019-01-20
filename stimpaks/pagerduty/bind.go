package pagerduty

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (p *Pagerduty) BindCommand(parentCommand *cobra.Command) {

	var cmd = &cobra.Command{
		Use:   "pagerduty",
		Short: "Send events to Pagerduty",
		Long:  `Sends trigger, acknowledge and resolve events to Pagerduty`,
		Run: func(cmd *cobra.Command, args []string) {
			p.log.Debug("Running `pagerduty` command")
			p.SendEvent()
		},
	}

	cmd.Flags().StringP("action", "a", "", "Required. The event action. Must be one of [trigger, acknowledge, resolve]")
	viper.BindPFlag("notify-pagerduty-action", cmd.Flags().Lookup("action"))

	cmd.Flags().StringP("service", "s", "", "Required. Name of Pagerduty service to send notification to")
	p.api.Config.BindPFlag("notify-pagerduty-service", cmd.Flags().Lookup("service"))

	cmd.Flags().StringP("severity", "r", "", "Required. How impacted the affected system is. Displayed to users in lists and influences the priority of any created incidents. Must be one of [Info, Warning, Error, Critical]")
	viper.BindPFlag("notify-pagerduty-severity", cmd.Flags().Lookup("severity"))

	cmd.Flags().StringP("summary", "m", "", "Required. A high-level, text summary message of the event. Will be used to construct an alert's description.")
	viper.BindPFlag("notify-pagerduty-summary", cmd.Flags().Lookup("summary"))

	cmd.Flags().StringP("source", "o", "", "Required. Specific human-readable unique identifier, such as a hostname, for the system having the problem.")
	viper.BindPFlag("notify-pagerduty-source", cmd.Flags().Lookup("source"))

	cmd.Flags().StringP("component", "c", "", "The part or component of the affected system that is broken.")
	viper.BindPFlag("notify-pagerduty-component", cmd.Flags().Lookup("component"))

	cmd.Flags().StringP("group", "g", "", "A cluster or grouping of sources. For example, sources “prod-datapipe-02” and “prod-datapipe-03” might both be part of “prod-datapipe”")
	viper.BindPFlag("notify-pagerduty-group", cmd.Flags().Lookup("group"))

	cmd.Flags().StringP("class", "l", "", "The class/type of the event")
	viper.BindPFlag("notify-pagerduty-class", cmd.Flags().Lookup("class"))

	cmd.Flags().StringP("details", "d", "", "Free-form details from the event")
	viper.BindPFlag("notify-pagerduty-details", cmd.Flags().Lookup("details"))

	cmd.Flags().StringP("dedupkey", "", "", "UniquedDe-duplication key for the alert. Should the same between all actions for a single incident")
	viper.BindPFlag("notify-pagerduty-dedupkey", cmd.Flags().Lookup("dedupkey"))

	viper.BindPFlag("vault-pagerduty-apikey-path", parentCommand.PersistentFlags().Lookup("vault-pagerduty-apikey-path"))
	viper.BindPFlag("vault-pagerduty-apikey-key", parentCommand.PersistentFlags().Lookup("vault-pagerduty-apikey-key"))

	parentCommand.AddCommand(cmd)
}

// type Event struct {
// 	Action    string `mapstructure:"notify-pagerduty-action"`
// 	Service   string `mapstructure:"notify-pagerduty-service"`
// 	Severity  string `mapstructure:"notify-pagerduty-severity"`
// 	Summary   string `mapstructure:"notify-pagerduty-summary"`
// 	Source    string `mapstructure:"notify-pagerduty-source"`
// 	Component string `mapstructure:"notify-pagerduty-component"`
// 	Group     string `mapstructure:"notify-pagerduty-group"`
// 	Class     string `mapstructure:"notify-pagerduty-class"`
// 	Details   string `mapstructure:"notify-pagerduty-details"`
// 	DedupKey  string `mapstructure:"notify-pagerduty-dedupkey"`
// }

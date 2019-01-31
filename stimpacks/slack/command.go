package slack

import (
	"github.com/readytalk/stim/stim"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type pepper struct {
	Name string
}

func (s *Slack) BindStim(stim *stim.Stim) {
	s.stim = stim
}

func (s *Slack) Command(viper *viper.Viper) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "slack",
		Short: "Interact with Slack",
		Long:  `Send/Recieve messages, etc. to/from Slack`,
		Run: func(cmd *cobra.Command, args []string) {
			s.postMessage()
		},
	}

	cmd.Flags().StringP("channel", "c", "", "Required. The channel name to send the message to")
	viper.BindPFlag("slack.channel", cmd.Flags().Lookup("channel"))

	cmd.Flags().StringP("message", "m", "", "Required. The message to send")
	viper.BindPFlag("slack.message", cmd.Flags().Lookup("message"))

	cmd.Flags().StringP("username", "u", "", "Username for the message to appear as")
	viper.BindPFlag("slack.username", cmd.Flags().Lookup("username"))
	viper.SetDefault("slack.username", "stim")

	cmd.Flags().StringP("icon-url", "i", "", "Url to use as the icon for the message")
	viper.BindPFlag("slack.icon-url", cmd.Flags().Lookup("icon-url"))
	viper.SetDefault("slack.icon-url", "https://vignette.wikia.nocookie.net/fallout/images/7/7e/FoS_stimpak.png/revision/latest")

	return cmd
}

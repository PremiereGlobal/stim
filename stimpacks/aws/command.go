package aws

import (
	"github.com/readytalk/stim/stim"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func (a *Aws) BindStim(stim *stim.Stim) {
	a.stim = stim
}

func (a *Aws) Command(viper *viper.Viper) *cobra.Command {

	var cmd = &cobra.Command{
		Use:   "aws",
		Short: "Interact with AWS",
		Long:  `Get credentials, etc.`,
		Run: func(cmd *cobra.Command, args []string) {
			a.GetCredentials()
		},
	}

	// cmd.Flags().StringP("channel", "c", "", "Required. The channel name to send the message to")
	// viper.BindPFlag("slack.channel", cmd.Flags().Lookup("channel"))
	//
	// cmd.Flags().StringP("message", "m", "", "Required. The message to send")
	// viper.BindPFlag("slack.message", cmd.Flags().Lookup("message"))
	//
	// cmd.Flags().StringP("username", "u", "", "Username for the message to appear as")
	// viper.BindPFlag("slack.username", cmd.Flags().Lookup("username"))
	//
	// cmd.Flags().StringP("icon-url", "i", "", "Url to use as the icon for the message")
	// viper.BindPFlag("slack.icon-url", cmd.Flags().Lookup("icon-url"))

	return cmd
}

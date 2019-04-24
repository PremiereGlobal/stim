package slack

import (
	"errors"
	slackpkg "github.com/PremiereGlobal/stim/pkg/slack"
)

func (s *Slack) postMessage() {

	var err error

	// Get a new authenticated Slack
	slack := s.stim.Slack()

	// Prompt for the channel name (if not provided)
	channelName := s.stim.GetConfig("slack.channel")
	if channelName == "" && s.stim.IsAutomated() {
		s.stim.Fatal(errors.New("Slack channel not specified"))
	} else if channelName == "" {

		// Get the channel list
		channels, err := slack.GetChannels()
		s.stim.Fatal(err)

		// Prompt for channel
		channelName, err = s.stim.PromptSearchList("Choose Channel:", channels)
		s.stim.Fatal(err)

	}

	// Prompt for the message (if not provided)
	text := s.stim.GetConfig("slack.message")
	if text == "" && s.stim.IsAutomated() {
		s.stim.Fatal(errors.New("Slack message not specified"))
	} else if text == "" {
		text, err = s.stim.PromptString("Message", "")
		s.stim.Fatal(err)
	}

	username := s.stim.GetConfig("slack.username")
	if username == "" && !s.stim.IsAutomated() {
		username, err = s.stim.PromptString("Display Name", DEFAULT_MESSAGE_USERNAME)
		s.stim.Fatal(err)
	}

	iconUrl := s.stim.GetConfig("slack.icon-url")
	if iconUrl == "" && !s.stim.IsAutomated() {
		iconUrl, err = s.stim.PromptString("Icon URL", DEFAULT_MESSAGE_ICON_URL)
		s.stim.Fatal(err)
	}

	// Construct the message
	message := &slackpkg.Message{
		Channel:  channelName,
		Username: username,
		Text:     text,
		IconUrl:  iconUrl,
	}

	// Post the message
	err = slack.PostMessage(message)
	s.stim.Fatal(err)

}

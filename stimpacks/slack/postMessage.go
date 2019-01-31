package slack

import (
	slackpkg "github.com/readytalk/stim/pkg/slack"
)

func (s *Slack) postMessage() {

	// Get a new authenticated Slack
	slack := s.stim.Slack()

	// Prompt for the channel name (if not provided)
	channelName := s.stim.GetConfig("slack.channel")
	if channelName == "" {

		// Get the channel list
		channels, err := slack.GetChannels()
		if err != nil {
			s.stim.Fatal(err)
			return
		}

		// Prompt for channel
		channelName, err = s.stim.PromptSearchList("Choose Channel:", channels, "")
		if err != nil {
			s.stim.Fatal(err)
			return
		}

	}

	// Prompt for the message (if not provided)
	text, err := s.stim.PromptString("Message", s.stim.GetConfig("slack.message"), "")
	if err != nil {
		s.stim.Fatal(err)
		return
	}

	// Construct the message
	message := &slackpkg.Message{
		Channel:  channelName,
		Username: s.stim.GetConfig("slack.username"),
		Text:     text,
		IconUrl:  s.stim.GetConfig("slack.icon-url"),
	}

	// Post the message
	err = slack.PostMessage(message)
	if err != nil {
		s.stim.Fatal(err)
		return
	}

}

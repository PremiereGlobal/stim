package slack

import (
	"errors"

	"github.com/nlopes/slack"
	"github.com/readytalk/stim/pkg/stimlog"
)

// Slack is the main object
type Slack struct {
	client *slack.Client
	log    *stimlog.StimLogger
	config *Config
}

// Config contains information about setting up a new slack client
type Config struct {
	Token string
}

// Message contains information about a slack message
type Message struct {
	Channel  string
	Username string
	Text     string
	IconUrl  string
}

// New builds a slack client from the provided config
func New(config *Config, sl *stimlog.StimLogger) (*Slack, error) {

	client := slack.New(config.Token)

	s := &Slack{config: config, client: client}
	if sl != nil {
		s.log = sl
	} else {
		s.log = stimlog.GetLogger()
	}

	return s, nil
}

// GetChannels looks up and returns a slice of channel names
func (s *Slack) GetChannels() ([]string, error) {
	channels, err := s.client.GetChannels(false)
	if err != nil {
		return nil, err
	}

	var result []string
	for _, channel := range channels {
		result = append(result, channel.Name)
	}

	return result, nil
}

func (s *Slack) getChannelIdByName(name string) (string, error) {
	channels, err := s.client.GetChannels(false)
	if err != nil {
		return "", err
	}

	for _, channel := range channels {
		if channel.Name == name {
			return channel.ID, nil
		}
	}

	return "", errors.New("Channel " + name + " not found")
}

// PostMessage posts a message to a Slack channel with the provided
// Message parameters
func (s *Slack) PostMessage(msg *Message) error {

	if msg.Text == "" {
		return errors.New("Slack message text required.")
	}

	id, err := s.getChannelIdByName(msg.Channel)
	if err != nil {
		return err
	}

	parameters := slack.NewPostMessageParameters()
	if msg.Username != "" {
		parameters.Username = msg.Username
	}
	if msg.IconUrl != "" {
		parameters.IconURL = msg.IconUrl
	}

	channelId, timestamp, err := s.client.PostMessage(id, slack.MsgOptionText(msg.Text, false), slack.MsgOptionPostMessageParameters(parameters))
	if err != nil {
		return err
	}

	s.log.Debug("Slack message successfully sent at " + timestamp + " to channel " + msg.Channel + " (" + channelId + ")")

	return nil
}

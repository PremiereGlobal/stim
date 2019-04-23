package slack

import (
	"errors"

	"github.com/PremiereGlobal/stim/pkg/stimlog"
	"github.com/nlopes/slack"
)

// Slack is the main object
type Slack struct {
	client *slack.Client
	log    Logger
	config *Config
}

// Config contains information about setting up a new slack client
type Config struct {
	Token string
	Log   Logger
}

// Message contains information about a slack message
type Message struct {
	Channel  string
	Username string
	Text     string
	IconUrl  string
}

type Logger interface {
	Debug(...interface{})
	Warn(...interface{})
	Fatal(...interface{})
}

// New builds a slack client from the provided config
func New(config *Config) (*Slack, error) {

	client := slack.New(config.Token)

	s := &Slack{config: config, client: client}
	if config.Log != nil {
		s.log = config.Log
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

package slack

import (
	"errors"
	"fmt"
	"github.com/nlopes/slack"
)

// Slack is the main object
type Slack struct {
	client *slack.Client
	config *Config
}

// Config contains information about setting up a new slack client
type Config struct {
	Token string
	Logger
}

// Message contains information about a slack message
type Message struct {
	Channel  string
	Username string
	Text     string
	IconUrl  string
}

// Logger is an interface for passing a custom logger to be used by this package
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
}

func (s *Slack) debug(message string) {
	if s.config.Logger != nil {
		s.config.Debug(message)
	}
}

func (s *Slack) info(message string) {
	if s.config.Logger != nil {
		s.config.Info(message)
	} else {
		fmt.Println(message)
	}
}

// New builds a slack client from the provided config
func New(config *Config) (*Slack, error) {

	client := slack.New(config.Token)

	s := &Slack{config: config, client: client}

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

	s.debug("Slack message successfully sent at " + timestamp + " to channel " + msg.Channel + " (" + channelId + ")")

	return nil
}

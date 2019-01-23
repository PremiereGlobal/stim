package prometheus

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"time"
)

type Prometheus struct {
	client  *api.Client
	config  *Config
	context context.Context
	API     v1.API
}

type Config struct {
	Address string
	Logger
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
}

func (p *Prometheus) Debug(message string) {
	if p.config.Logger != nil {
		p.config.Debug(message)
	}
}

func (p *Prometheus) Info(message string) {
	if p.config.Logger != nil {
		p.config.Info(message)
	} else {
		fmt.Println(message)
	}
}

func New(config *Config) (*Prometheus, error) {

	apiConfig := api.Config{Address: config.Address}
	client, err := api.NewClient(apiConfig)
	if err != nil {
		return nil, err
	}

	api := v1.NewAPI(client)

	p := &Prometheus{client: &client, API: api, context: context.Background()}

	return p, nil
}

func (p *Prometheus) QueryInstant(query string) (string, error) {
	result, err := p.API.Query(p.context, query, time.Now())
	if err != nil {
		return "", err
	}

	t := result.Type()
	var d []byte
	t.UnmarshalJSON(d)

	return string(d), nil
}

package prometheus

import (
	"context"
	"time"

	"github.com/PremiereGlobal/stim/pkg/stimlog"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
)

type Prometheus struct {
	client  *api.Client
	config  *Config
	context context.Context
	API     v1.API
	log     Logger
}

type Config struct {
	Address string
	Log     Logger
}

type Logger interface {
	Debug(...interface{})
	Warn(...interface{})
	Fatal(...interface{})
}

func New(config *Config) (*Prometheus, error) {

	apiConfig := api.Config{Address: config.Address}
	client, err := api.NewClient(apiConfig)
	if err != nil {
		return nil, err
	}

	api := v1.NewAPI(client)

	p := &Prometheus{client: &client, API: api, context: context.Background()}
	if config.Log != nil {
		p.log = config.Log
	} else {
		p.log = stimlog.GetLogger()
	}
	return p, nil
}

func (p *Prometheus) QueryInstant(query string) (string, error) {
	result, warnings, err := p.API.Query(p.context, query, time.Now())
	if err != nil {
		return "", err
	}
  p.log.Warn(warnings)

	t := result.Type()
	var d []byte
	t.UnmarshalJSON(d)

	return string(d), nil
}

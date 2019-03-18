package prometheus

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/readytalk/stim/pkg/stimlog"
)

type Prometheus struct {
	client  *api.Client
	config  *Config
	context context.Context
	API     v1.API
	log     *stimlog.StimLogger
}

type Config struct {
	Address string
}

func New(config *Config, sl *stimlog.StimLogger) (*Prometheus, error) {

	apiConfig := api.Config{Address: config.Address}
	client, err := api.NewClient(apiConfig)
	if err != nil {
		return nil, err
	}

	api := v1.NewAPI(client)

	p := &Prometheus{client: &client, API: api, context: context.Background()}
	if sl != nil {
		p.log = sl
	} else {
		p.log = stimlog.GetLogger()
	}
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

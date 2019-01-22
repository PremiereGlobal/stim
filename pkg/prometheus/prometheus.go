package prometheus

import (
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	// "time"
	"context"
)

type Prometheus struct {
	client  *api.Client
	config  *Config
	context context.Context
	API     v1.API
}

type Config struct {
	Address string
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

// func (p *Prometheus) QueryInstant(query string) (string, error) {
// 	result, err := p.client.API.Query(p.context, query, time.Now())
// 	if err != nil {
// 		return "", err
// 	}
//
// 	t := result.Type()
// 	var d []byte
// 	t.UnmarshalJSON(d)
//
// 	return string(d), nil
// }

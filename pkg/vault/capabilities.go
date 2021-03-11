package vault

import (
	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/mapstructure"

	"context"
)

type CapabilitiesSelfOptions struct {
	Paths []string `json:"paths,omitempty"`
}

type CapabilitiesSelfResults map[string][]string

func (v *Vault) CapabilitiesSelf(opts *CapabilitiesSelfOptions) (CapabilitiesSelfResults, error) {
	request := v.client.NewRequest("POST", "/v1/sys/capabilities-self")
	if err := request.SetJSONBody(opts); err != nil {
		return nil, err
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	response, err := v.client.RawRequestWithContext(ctx, request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	parsedSecret, err := api.ParseSecret(response.Body)
	if err != nil {
		return nil, err
	}

	var results CapabilitiesSelfResults
	err = mapstructure.Decode(parsedSecret.Data, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

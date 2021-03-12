package vault

import (
	"context"
	"encoding/json"
)

type CapabilitiesSelfOptions struct {
	Paths []string `json:"paths,omitempty"`
}

type CapabilitiesSelfResults struct {
	Data map[string][]string `json:"data"`
}

func (v *Vault) CapabilitiesSelf(opts *CapabilitiesSelfOptions) (*CapabilitiesSelfResults, error) {
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

	results := CapabilitiesSelfResults{}
	err = json.NewDecoder(response.Body).Decode(&results)
	if err != nil {
		return nil, err
	}

	return &results, nil
}

package vault

import (
	"encoding/json"
	"time"
)

// GetCurrentTokenTTL gets the TTL of the current token
func (v *Vault) GetCurrentTokenTTL() (time.Duration, error) {

	// Get the token info from Vault
	secret, err := v.client.Auth().Token().LookupSelf()
	if err != nil {
		return 0, err
	}

	// Get our TTL from the Vault secret interface{}
	ttl, err := secret.Data["ttl"].(json.Number).Int64()
	if err != nil {
		return 0, err
	}

	// Convert our ttl (int64) to time.Duration
	duration := time.Duration(ttl) * time.Second

	return duration, nil
}

package vault

import (
	"time"
)

// Renew lease takes a Vault lease ID and renews it for the provided duration
// Returns the actual renew time (may be different than requested)
func (v *Vault) RenewLease(leaseID string, duration time.Duration) (time.Duration, error) {

	v.log.Debug("Renewing lease " + leaseID + " for " + duration.String())
	secret, err := v.client.Sys().Renew(leaseID, int(duration.Seconds()))
	if err != nil {
		return 0, err
	}

	leaseDuration := time.Duration(secret.LeaseDuration) * time.Second

	return leaseDuration, nil
}

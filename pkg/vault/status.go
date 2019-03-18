package vault

import (
	"strconv"
	"time"
)

func (v *Vault) isVaultHealthy() (bool, error) {
	result, err := v.client.Sys().Health()
	if err != nil {
		return false, err
	}

	v.log.Debug("Vault server info from (" + v.client.Address() + ")")
	v.log.Debug("  Initialized: " + strconv.FormatBool(result.Initialized))
	v.log.Debug("  Sealed: " + strconv.FormatBool(result.Sealed))
	v.log.Debug("  Standby: " + strconv.FormatBool(result.Standby))
	v.log.Debug("  Version: " + result.Version)
	v.log.Debug("  ClusterName: " + result.ClusterName)
	v.log.Debug("  ClusterID: " + result.ClusterID)
	v.log.Debug("  ServerTime: (" + strconv.FormatInt(result.ServerTimeUTC, 10) + ") " + time.Unix(result.ServerTimeUTC, 0).UTC().String())

	return true, nil
}

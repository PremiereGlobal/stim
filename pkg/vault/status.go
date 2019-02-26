package vault

import (
	"strconv"
	"time"
)

func (v *Vault) isVaultHealthy() (bool, error) {

	result, err := v.client.Sys().Health()
	if err != nil {
		return false, v.parseError(err)
	}

	v.Debug("Vault server info from (" + v.client.Address() + ")")
	v.Debug("  Initialized: " + strconv.FormatBool(result.Initialized))
	v.Debug("  Sealed: " + strconv.FormatBool(result.Sealed))
	v.Debug("  Standby: " + strconv.FormatBool(result.Standby))
	v.Debug("  Version: " + result.Version)
	v.Debug("  ClusterName: " + result.ClusterName)
	v.Debug("  ClusterID: " + result.ClusterID)
	v.Debug("  ServerTime: (" + strconv.FormatInt(result.ServerTimeUTC, 10) + ") " + time.Unix(result.ServerTimeUTC, 0).UTC().String())

	return true, nil
}

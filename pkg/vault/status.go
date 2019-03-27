package vault

import (
	"github.com/readytalk/stim/pkg/log"

	"strconv"
	"time"
)

func (v *Vault) isVaultHealthy() (bool, error) {

	result, err := v.client.Sys().Health()
	if err != nil {
		return false, v.parseError(err)
	}

	log.Debug("Vault server info from (" + v.client.Address() + ")")
	log.Debug("  Initialized: " + strconv.FormatBool(result.Initialized))
	log.Debug("  Sealed: " + strconv.FormatBool(result.Sealed))
	log.Debug("  Standby: " + strconv.FormatBool(result.Standby))
	log.Debug("  Version: " + result.Version)
	log.Debug("  ClusterName: " + result.ClusterName)
	log.Debug("  ClusterID: " + result.ClusterID)
	log.Debug("  ServerTime: (" + strconv.FormatInt(result.ServerTimeUTC, 10) + ") " + time.Unix(result.ServerTimeUTC, 0).UTC().String())

	return true, nil
}

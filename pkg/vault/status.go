package vault

func (v *Vault) isVaultHealthy() (bool, error) {
	_, err := v.client.Sys().Health()
	if err != nil {
		return false, err
	}

	// log.Debug("Vault server info from (", v.client.Address(), ")")
	// log.Debug("  Initialized: ", result.Initialized)
	// log.Debug("  Sealed: ", result.Sealed)
	// log.Debug("  Standby: ", result.Standby)
	// log.Debug("  Version: ", result.Version)
	// log.Debug("  ClusterName: ", result.ClusterName)
	// log.Debug("  ClusterID: ", result.ClusterID)
	// log.Debug("  ServerTime: (", result.ServerTimeUTC, ") ", time.Unix(result.ServerTimeUTC, 0).UTC())
	// log.Debug("  Standby: ", result.Standby)

	return true, nil
}

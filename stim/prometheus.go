package stim

import (
	"github.com/PremiereGlobal/stim/pkg/prometheus"
)

func (stim *Stim) Prometheus() *prometheus.Prometheus {
	stim.log.Debug("Stim-Prometheus: Creating")

	address := stim.ConfigGetString("prometheus.address")
	stim.log.Debug("Stim-Prometheus: Using Address {}", address)

	p, err := prometheus.New(&prometheus.Config{Address: address, Log: stim.log})
	if err != nil {
		stim.log.Fatal("Stim-Prometheus: Error Initializaing: {}", err)
	}

	return p
}

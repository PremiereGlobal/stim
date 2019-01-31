package discover

import (
	"github.com/readytalk/stim/stim"
)

type Discover struct {
	name string
	stim *stim.Stim
}

func New() *Discover {
	discover := &Discover{name: "discover"}
	return discover
}

func (d *Discover) Name() string {
	return d.name
}

func (d *Discover) DiscoverClusters() (string, error) {
	p := d.stim.Prometheus()
	result, err := p.QueryInstant("kubernetes_build_info{}")
	if err != nil {
		return "", err
	}

	return result, nil
}

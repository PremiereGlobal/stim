package deploy

import (
	"errors"
	"fmt"
	"github.com/PremiereGlobal/stim/stim"
	"github.com/davecgh/go-spew/spew"
)

type Deploy struct {
	name   string
	stim   *stim.Stim
	config Config
}

func New() *Deploy {
	deploy := &Deploy{}
	return deploy
}

func (d *Deploy) Name() string {
	return d.name
}

func (d *Deploy) Run() {

	d.ParseConfig()

	environments := make([]string, 0, len(d.config.Environments))
	environment := ""
	environmentArg := d.stim.GetConfig("deploy.environment")
	if environmentArg != "" {
		if _, ok := d.config.Environments[environmentArg]; ok {
			environment = environmentArg
		} else {
			d.stim.Fatal(errors.New(fmt.Sprintf("Provided environment value '%s' is not in config file", environmentArg)))
		}
	} else {
		for e := range d.config.Environments {
			environments = append(environments, e)
		}
		environment, _ = d.stim.PromptList("Which environment?", environments, d.stim.GetConfig("deploy.environment"))
	}

	clusters := make([]string, 0, len(d.config.Environments[environment].Clusters)+1)
	cluster := ""
	clusterArg := d.stim.GetConfig("deploy.cluster")
	if clusterArg != "" {
		if _, ok := d.config.Environments[environment].Clusters[clusterArg]; ok {
			cluster = clusterArg
		} else {
			d.stim.Fatal(errors.New(fmt.Sprintf("Provided cluster value '%s' is not in config file under environment '%s'", clusterArg, environment)))
		}
	} else {
		clusters := append(clusters, "--ALL--")
		for c := range d.config.Environments[environment].Clusters {
			clusters = append(clusters, c)
		}
		cluster, _ = d.stim.PromptList("Which cluster?", clusters, d.stim.GetConfig("deploy.cluster"))
	}

	spew.Dump(fmt.Sprintf("Deploying to %s, %s", environment, cluster))
	spew.Dump(d.config)
}

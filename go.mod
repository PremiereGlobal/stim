module github.com/PremiereGlobal/stim

go 1.13

replace (
	github.com/docker/docker => github.com/docker/engine v0.0.0-20190822205725-ed20165a37b4
	k8s.io/api => k8s.io/api v0.0.0-20190222213804-5cb15d344471
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190221213512-86fb29eff628
	k8s.io/client-go => k8s.io/client-go v10.0.0+incompatible
)

require (
	github.com/PagerDuty/go-pagerduty v0.0.0-20191002190746-f60f4fc45222
	github.com/PremiereGlobal/vault-to-envs v0.2.1
	github.com/aws/aws-sdk-go v1.25.6
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/cornelk/hashmap v1.0.0
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/go-ini/ini v1.48.0
	github.com/hashicorp/vault v1.2.3
	github.com/hashicorp/vault/api v1.0.5-0.20190909201928-35325e2c3262
	github.com/imdario/mergo v0.3.8
	github.com/manifoldco/promptui v0.3.2
	github.com/mitchellh/go-homedir v1.1.0
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/nlopes/slack v0.6.0
	github.com/prometheus/client_golang v1.1.0
	github.com/skratchdot/open-golang v0.0.0-20190402232053-79abb63cd66e
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	golang.org/x/crypto v0.0.0-20191002192127-34f69633bfdc
	gopkg.in/alecthomas/kingpin.v3-unstable v3.0.0-20180810215634-df19058c872c // indirect
	gopkg.in/yaml.v2 v2.2.4
	gopkg.in/yaml.v3 v3.0.0-20190924164351-c8b7dadae555
	gotest.tools v2.2.0+incompatible
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/utils v0.0.0-20190923111123-69764acb6e8e // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)

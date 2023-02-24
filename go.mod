module github.com/PremiereGlobal/stim

go 1.13

replace (
	github.com/PremiereGlobal/stim => ./
	github.com/docker/docker => github.com/docker/engine v0.0.0-20190822205725-ed20165a37b4
	k8s.io/api => k8s.io/api v0.0.0-20190222213804-5cb15d344471
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190221213512-86fb29eff628
	k8s.io/client-go => k8s.io/client-go v10.0.0+incompatible
)

require (
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/PagerDuty/go-pagerduty v0.0.0-20191002190746-f60f4fc45222
	github.com/PremiereGlobal/vault-to-envs v0.2.2-0.20190928170516-b94151c229ae
	github.com/aws/aws-sdk-go v1.25.6
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/cornelk/hashmap v1.0.0
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/go-ini/ini v1.48.0
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/hashicorp/vault v1.2.3
	github.com/hashicorp/vault/api v1.0.5-0.20190909201928-35325e2c3262
	//	github.com/hashicorp/vault v1.0.2
	github.com/imdario/mergo v0.3.8
	github.com/krolaw/zipstream v0.0.0-20180621105154-0a2661891f94
	github.com/manifoldco/promptui v0.3.2
	github.com/mitchellh/go-homedir v1.1.0
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/nicksnyder/go-i18n v1.10.0 // indirect
	github.com/nlopes/slack v0.6.0
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/prometheus/client_golang v1.1.0
	github.com/skratchdot/open-golang v0.0.0-20190402232053-79abb63cd66e
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/tevino/abool v1.2.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/lint v0.0.0-20191125180803-fdd1cda4f05f // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4
	golang.org/x/net v0.7.0 // indirect
	gopkg.in/alecthomas/kingpin.v3-unstable v3.0.0-20180810215634-df19058c872c // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.8
	gopkg.in/yaml.v3 v3.0.0-20190924164351-c8b7dadae555
	gotest.tools v2.2.0+incompatible
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v0.3.0 // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)

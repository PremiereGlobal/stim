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
	4d63.com/gochecknoglobals v0.0.0-20190306162314-7c3491d2b6ec // indirect
	4d63.com/gochecknoinits v0.0.0-20200108094044-eb73b47b9fc4 // indirect
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/PagerDuty/go-pagerduty v0.0.0-20191002190746-f60f4fc45222
	github.com/PremiereGlobal/vault-to-envs v0.2.2-0.20190928170516-b94151c229ae
	github.com/alecthomas/gocyclo v0.0.0-20150208221726-aa8f8b160214 // indirect
	github.com/alexkohler/nakedret v1.0.0 // indirect
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
	github.com/jgautheron/goconst v0.0.0-20170703170152-9740945f5dcb // indirect
	github.com/krolaw/zipstream v0.0.0-20180621105154-0a2661891f94
	github.com/manifoldco/promptui v0.3.2
	github.com/mdempsky/maligned v0.0.0-20180708014732-6e39bd26a8c8 // indirect
	github.com/mdempsky/unconvert v0.0.0-20190921185256-3ecd357795af // indirect
	github.com/mibk/dupl v1.0.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/nicksnyder/go-i18n v1.10.0 // indirect
	github.com/nlopes/slack v0.6.0
	github.com/opennota/check v0.0.0-20180911053232-0c771f5545ff // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/prometheus/client_golang v1.1.0
	github.com/securego/gosec v0.0.0-20200203094520-d13bb6d2420c // indirect
	github.com/skratchdot/open-golang v0.0.0-20190402232053-79abb63cd66e
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/stripe/safesql v0.2.0 // indirect
	github.com/walle/lll v1.0.1 // indirect
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
	gopkg.in/alecthomas/kingpin.v3-unstable v3.0.0-20180810215634-df19058c872c // indirect
	gopkg.in/yaml.v2 v2.2.8
	gopkg.in/yaml.v3 v3.0.0-20190924164351-c8b7dadae555
	gotest.tools v2.2.0+incompatible
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v0.3.0 // indirect
	mvdan.cc/interfacer v0.0.0-20180901003855-c20040233aed // indirect
	mvdan.cc/lint v0.0.0-20170908181259-adc824a0674b // indirect
	mvdan.cc/unparam v0.0.0-20191111180625-960b1ec0f2c2 // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)

module github.com/PremiereGlobal/stim

go 1.12

replace (
	git.apache.org/thrift.git => github.com/apache/thrift v0.12.0
	github.com/PremiereGlobal/stim => ./
	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.4.2
	github.com/docker/docker => github.com/docker/engine v0.0.0-20190822205725-ed20165a37b4
	github.com/hashicorp/vault/api => github.com/hashicorp/vault/api v1.0.2
	k8s.io/api => k8s.io/api v0.0.0-20190222213804-5cb15d344471
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190221213512-86fb29eff628
	k8s.io/client-go => k8s.io/client-go v10.0.0+incompatible
)

require (
	github.com/PagerDuty/go-pagerduty v0.0.0-20181104233218-fe8f9c4593d0
	github.com/PremiereGlobal/vault-to-envs v0.2.1-0.20190908001744-dd472f3fd11e
	github.com/aws/aws-sdk-go v1.20.20
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/cornelk/hashmap v1.0.0
	github.com/docker/distribution v2.6.0-rc.1.0.20180327202408-83389a148052+incompatible // indirect
	github.com/docker/docker v0.0.0-00010101000000-000000000000
	github.com/go-ini/ini v1.42.0
	github.com/google/gofuzz v1.0.0 // indirect
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/hashicorp/vault v1.2.2
	github.com/hashicorp/vault/api v1.0.5-0.20190814205542-3b036e58e950
	// github.com/hashicorp/vault v1.0.2
	github.com/imdario/mergo v0.3.7 // indirect
	github.com/lusis/go-slackbot v0.0.0-20180109053408-401027ccfef5 // indirect
	github.com/lusis/slack-test v0.0.0-20190426140909-c40012f20018 // indirect
	github.com/manifoldco/promptui v0.3.2
	github.com/mitchellh/go-homedir v1.1.0
	github.com/morikuni/aec v0.0.0-20170113033406-39771216ff4c // indirect
	github.com/nicksnyder/go-i18n v1.10.0 // indirect
	github.com/nlopes/slack v0.5.0
	github.com/prometheus/client_golang v0.9.3
	github.com/skratchdot/open-golang v0.0.0-20190104022628-a2dfa6d0dab6
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	golang.org/x/crypto v0.0.0-20190829043050-9756ffdc2472
	golang.org/x/lint v0.0.0-20190409202823-959b441ac422 // indirect
	golang.org/x/net v0.0.0-20190827160401-ba9fcec4b297 // indirect
	golang.org/x/sys v0.0.0-20190904154756-749cb33beabd // indirect
	golang.org/x/tools v0.0.0-20190904213738-958971f5c2bf // indirect
	gopkg.in/alecthomas/kingpin.v3-unstable v3.0.0-20180810215634-df19058c872c // indirect
	gopkg.in/yaml.v2 v2.2.2
	k8s.io/client-go v0.0.0-20190419212732-59781b88d0fa
	k8s.io/klog v0.3.0 // indirect
	sigs.k8s.io/yaml v1.1.0 // indirect
)

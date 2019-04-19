package stim

import (
	"github.com/readytalk/stim/pkg/aws"
)

func (stim *Stim) Aws(accessKey string, secretKey string) *aws.Aws {
	stim.GetLogger().Debug("Stim-Aws: Creating")
	a, err := aws.New(&aws.Config{AccessKey: accessKey, SecretKey: secretKey, Log: stim.GetLogger()})
	if err != nil {
		stim.log.Fatal("Stim-Aws: Error Initializaing: ", err)
	}

	return a
}

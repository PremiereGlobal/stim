package stim

import (
	"github.com/readytalk/stim/pkg/aws"
)

func (stim *Stim) Aws() *aws.Aws {
	stim.log.Debug("Stim-Aws: Creating")
	// vault = stim.Vault()
	// token, err := vault.GetSecretKey("secret/slack/stimbot", "apikey")
	// if err != nil {
	// 	stim.log.Fatal(err)
	// }

	a, err := aws.New(&aws.Config{Log: stim.log})
	if err != nil {
		stim.log.Fatal("Stim-Aws: Error Initializaing: ", err)
	}

	return a
}

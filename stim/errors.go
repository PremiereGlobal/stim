package stim

func (stim *Stim) Debug(err error) {
	if err != nil {
		stim.log.Debug(err)
	}
}

func (stim *Stim) Fatal(err error) {
	if err != nil {
		stim.log.Fatal(err)
	}
}

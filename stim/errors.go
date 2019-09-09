package stim

func (stim *Stim) Debug(message string) {
	if message != "" {
		stim.log.Debug(message)
	}
}

func (stim *Stim) Warn(message string) {
	if message != "" {
		stim.log.Warn(message)
	}
}

func (stim *Stim) DebugError(err error) {
	if err != nil {
		stim.log.Debug(err)
	}
}

func (stim *Stim) Fatal(err error) {
	if err != nil {
		stim.log.Fatal(err)
	}
}

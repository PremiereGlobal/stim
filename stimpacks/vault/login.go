package vault

func (v *Vault) Login() {
	// Get a new Vault from the API
	_ = v.stim.Vault()

}

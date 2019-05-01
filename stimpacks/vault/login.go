package vault

// Login will connect to Vault server and login
func (v *Vault) Login() {
	// Get a new Vault from the API
	_ = v.stim.Vault()
}

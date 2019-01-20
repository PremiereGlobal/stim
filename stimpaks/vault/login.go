package vault

func (v *Vault) Login() {

	// Get a new Vault from the API
	stimVault := v.api.Vault()

	// Login
	stimVault.Login()
}

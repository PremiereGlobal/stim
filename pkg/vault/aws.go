package vault

import (
   "github.com/skratchdot/open-golang/open"
)

func (v *Vault) AWS() error {

  err := open.Run("https://google.com/")
  if err != nil {
		return err
	}

  return nil
}

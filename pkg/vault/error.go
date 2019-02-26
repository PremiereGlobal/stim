package vault

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"syscall"
)

// VaultError is the custom error type for this package
type VaultError struct {
	MessageParts  []string
	OriginalError error
}

// Error returns the error string
func (verr VaultError) Error() string {
	return fmt.Sprintf("Vault Error: %s", strings.Join(verr.MessageParts, "; "))
}

// parseError parses known errors into more user-friendly messages
func (v *Vault) parseError(err error) VaultError {

	var verr VaultError
	verr.OriginalError = err

	// Catch some known HTTP errors
	if uerr, ok := err.(*url.Error); ok {
		if oerr, ok := uerr.Err.(*net.OpError); ok {
			if addr, ok := oerr.Addr.(*net.TCPAddr); ok {
				if addr.IP.String() == "127.0.0.1" {
					verr.MessageParts = append(verr.MessageParts, "Vault appears to be connecting to localhost, ensure correct Vault address is set")
				}
			}

			if serr, ok := oerr.Err.(*os.SyscallError); ok {
				if serr.Err == syscall.ECONNREFUSED {
					verr.MessageParts = append(verr.MessageParts, "Connection Refused")
				}
			}
		}
	}

	if err == context.DeadlineExceeded {
		verr.MessageParts = append(verr.MessageParts, fmt.Sprintf("Timeout connecting after %v seconds. Ensure connectivity to Vault.", v.config.Timeout))
	}

	verr.MessageParts = append(verr.MessageParts, fmt.Sprintf("%v", err))

	return verr
}

// newError returns a new error based on a given string
func (v *Vault) newError(msg string) VaultError {
	return v.parseError(errors.New(msg))
}

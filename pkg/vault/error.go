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

// CustomVaultError is the custom error type for this package
type CustomVaultError struct {
	MessageParts  []string
	OriginalError error
}

// customError is custom error interface with one method that will block other functions
// from using this Error. This is not interchangeable with the standard error.
type CustomError interface {
	Error() string
	blockInterface()
}

// Error returns the error string
func (verr CustomVaultError) Error() string {
	return fmt.Sprintf("Vault: %s", strings.Join(verr.MessageParts, "; "))
}

//
func (verr CustomVaultError) blockInterface() {
}

// parseError parses known errors into more user-friendly messages
func (v *Vault) parseError(err error) CustomError {

	// Provent parseError from calling parseError again
	if serr, ok := err.(CustomError); ok {
		return serr
	}

	verr := &CustomVaultError{OriginalError: err}

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
func (v *Vault) newError(msg string) CustomError {
	return v.parseError(errors.New(msg))
}

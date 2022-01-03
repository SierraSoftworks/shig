package publickeys

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

var (
	ErrUntrustedNamespace = fmt.Errorf("ssign: the signature is using an untrusted namespace")
	ErrUntrustedKey       = fmt.Errorf("ssign: the signature is using an untrusted key")
)

// A Validator is responsible for checking that a given public key
// is trusted to generate signatures for a given namespace.
type Validator interface {
	Validate(key ssh.PublicKey) error
}

package publickeys

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

func NewSshKeyValidator(publicKey string) (Validator, error) {
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(publicKey))
	if err != nil {
		return nil, errors.Wrap(err, "shig: failed to parse public key")
	}

	return &allowedKeysValidator{
		keys: []ssh.PublicKey{pubKey},
	}, nil
}

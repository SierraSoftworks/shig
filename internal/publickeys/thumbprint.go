package publickeys

import "golang.org/x/crypto/ssh"

func NewThumbprintValidator(thumbprint string) Validator {
	return &thumbprintValidator{
		Thumbprint: thumbprint,
	}
}

type thumbprintValidator struct {
	Thumbprint string
}

func (v *thumbprintValidator) Validate(key ssh.PublicKey) error {
	if v.Thumbprint != ssh.FingerprintSHA256(key) {
		return ErrUntrustedKey
	}

	return nil
}

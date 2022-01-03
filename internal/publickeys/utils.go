package publickeys

import "golang.org/x/crypto/ssh"

type allowedKeysValidator struct {
	keys []ssh.PublicKey
}

func (v *allowedKeysValidator) Validate(key ssh.PublicKey) error {
	for _, pubKey := range v.keys {
		if pubKey.Type() != key.Type() {
			continue
		}

		if keysEqual(pubKey, key) {
			return nil
		}
	}

	return ErrUntrustedKey
}

func keysEqual(a, b ssh.PublicKey) bool {
	ab, bb := a.Marshal(), b.Marshal()

	if len(ab) != len(bb) {
		return false
	}

	for i := range ab {
		if ab[i] != bb[i] {
			return false
		}
	}

	return true
}

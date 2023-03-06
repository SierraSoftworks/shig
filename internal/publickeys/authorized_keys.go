package publickeys

import (
	"bufio"
	"os"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

func NewAuthorizedKeysFileValidator(path string) (Validator, error) {
	v := &allowedKeysValidator{
		keys: make([]ssh.PublicKey, 0),
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "shig: failed to open authorized_keys file")
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		pubKey, _, _, _, err := ssh.ParseAuthorizedKey(scanner.Bytes())
		if err != nil {
			return nil, err
		}

		v.keys = append(v.keys, pubKey)
	}

	return v, nil
}

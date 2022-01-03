package publickeys

import (
	"bufio"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

func NewGitHubValidator(username string) (Validator, error) {
	v := &allowedKeysValidator{
		keys: make([]ssh.PublicKey, 0),
	}

	resp, err := http.Get(fmt.Sprintf("https://github.com/%s.keys", username))
	if err != nil {
		return nil, errors.Wrap(err, "ssign: failed to fetch GitHub public keys")
	}

	scanner := bufio.NewScanner(resp.Body)
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

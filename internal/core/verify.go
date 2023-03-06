package core

import (
	"os"
	"strings"

	"github.com/SierraSoftworks/shig/internal/publickeys"
	"github.com/SierraSoftworks/sshsign-go"
	"golang.org/x/crypto/ssh"
)

type Verifier struct {
	output    Output
	sigFile   string
	verifier  sshsign.Verifier
	validator publickeys.Validator
}

func NewVerifier(
	output Output,
	namespace, hash, sigFile string,
	validator publickeys.Validator,
) *Verifier {
	verifier := sshsign.DefaultVerifier(namespace, hash)

	return &Verifier{
		output:    output,
		sigFile:   sigFile,
		verifier:  verifier,
		validator: validator,
	}
}

func (v *Verifier) Verify(file string) error {
	f, err := os.Open(file)
	if err != nil {
		v.output.Printf("FAIL: '%s' could not be opened\n", file)
		return err
	}
	defer f.Close()

	sigFile := strings.ReplaceAll(v.sigFile, "%f", file)
	sf, err := os.ReadFile(sigFile)
	if err != nil {
		v.output.Printf("FAIL: '%s' does not have a corresponding signature file '%s'\n", file, sigFile)
		return err
	}

	sig, _, err := sshsign.UnmarshalArmoured(sf)
	if err != nil {
		v.output.Printf("FAIL: '%s' is not a well-formatted signature file.\n", sigFile)
		return err
	}

	if err := v.verifier.Verify(f, sig); err != nil {
		v.output.Printf("FAIL: '%s' does not match the signature file '%s'\n", file, sigFile)
		return err
	}

	key, err := sig.GetPublicKey()
	if err != nil {
		v.output.Printf("FAIL: '%s' does not contain a valid public key in its signature\n", file)
		return err
	}

	if v.validator != nil {
		if err := v.validator.Validate(key); err != nil {
			v.output.Printf("FAIL: '%s' is signed by an untrusted key: %s\n", file, ssh.FingerprintSHA256(key))
			return err
		}
	}

	v.output.Printf("PASS: '%s' is signed by '%s'\n", file, ssh.FingerprintSHA256(key))
	return nil
}

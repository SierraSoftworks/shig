package core

import (
	"os"
	"strings"

	"github.com/SierraSoftworks/sshsign-go"
	"golang.org/x/crypto/ssh"
)

type Signer struct {
	output  Output
	sigFile string
	signer  sshsign.Signer
}

func NewSigner(
	output Output,
	key, namespace, hash, sigFile string,
) (*Signer, error) {
	pkc, err := os.ReadFile(os.ExpandEnv(key))
	if err != nil {
		output.Println("FAIL: Unable to read your SSH private key. Make sure that you have entered its path correctly and have permission to access it.")
		return nil, err
	}

	pk, err := ssh.ParsePrivateKey(pkc)
	if err != nil {
		output.Println("FAIL: Unable to parse your SSH private key. Make sure that it is a well-formatted SSH private key file.")
		return nil, err
	}

	signer := sshsign.DefaultSigner(namespace, hash, pk)

	return &Signer{
		sigFile: sigFile,
		output:  output,
		signer:  signer,
	}, nil
}

func (s *Signer) Sign(file string) error {
	f, err := os.Open(file)
	if err != nil {
		s.output.Printf("FAIL: '%s' could not be opened for signing.\n", file)
		return err
	}
	defer f.Close()

	sig, err := s.signer.Sign(f)
	if err != nil {
		s.output.Printf("FAIL: '%s' could not be signed.\n", file)
		return err
	}

	armoured, err := sig.MarshalArmoured()
	if err != nil {
		s.output.Printf("FAIL: '%s' could not format the signature file correctly.\n", file)
		return err
	}

	sigFile := strings.ReplaceAll(s.sigFile, "%f", file)
	if err := os.WriteFile(sigFile, armoured, 0644); err != nil {
		s.output.Printf("FAIL: '%s' could not be saved.\n", sigFile)
		return err
	}

	s.output.Printf("PASS: '%s' has been signed.\n", file)
	return nil
}

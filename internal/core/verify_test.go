package core_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path"
	"testing"

	"github.com/SierraSoftworks/shig/internal/core"
	"github.com/SierraSoftworks/shig/internal/publickeys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ssh"
)

func TestVerifier(t *testing.T) {
	home := TempDir(t)

	pk, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	pkf, err := os.OpenFile(path.Join(home, "id_rsa"), os.O_CREATE|os.O_WRONLY, 0600)
	require.NoError(t, err)

	require.NoError(t, pem.Encode(pkf, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(pk),
	}), "failed to write private key")
	pkf.Close()

	require.NoError(t, os.WriteFile(path.Join(home, "example.txt"), []byte("Hello, World!"), 0600))

	signer, err := core.NewSigner(&testOutput{t}, path.Join(home, "id_rsa"), "test", "sha512", "%f.sig")
	require.NoError(t, err)

	assert.NoError(t, signer.Sign(path.Join(home, "example.txt")))

	pubk, err := ssh.NewPublicKey(&pk.PublicKey)
	require.NoError(t, err)

	pubks := string(ssh.MarshalAuthorizedKey(pubk))
	t.Log("Public Key:", pubks)

	val, err := publickeys.NewSshKeyValidator(pubks)
	require.NoError(t, err, "failed to create validator")

	verifier := core.NewVerifier(&testOutput{t}, "test", "sha512", "%f.sig", val)
	assert.NoError(t, verifier.Verify(path.Join(home, "example.txt")))

	require.NoError(t, os.WriteFile(path.Join(home, "example.txt"), []byte("Hello, Claire!"), 0600))
	assert.Error(t, verifier.Verify(path.Join(home, "example.txt")), "signature should not match after file is modified")
}

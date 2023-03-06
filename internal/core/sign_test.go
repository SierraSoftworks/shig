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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSigner(t *testing.T) {
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

	sig, err := os.Stat(path.Join(home, "example.txt.sig"))
	assert.NoError(t, err, "failed to stat signature file")
	assert.Greater(t, sig.Size(), int64(0), "signature file is empty")
}

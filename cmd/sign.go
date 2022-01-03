/*
Copyright © 2022 Benjamin Pannell <admin@sierrasoftworks.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/SierraSoftworks/sshsign-go"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign a file using your SSH key.",
	Long: `Generate a cryptographic signature for a file using your SSH private key,
which can then be verified by anyone who knows your SSH public key. If you have your
key registered with GitHub, people can verify the signature using your GitHub username.`,
	Aliases: []string{"s"},
	Example: `ssign sign file1.txt file2.txt`,
	Annotations: map[string]string{
		"Group": "Signing",
	},
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key, _ := cmd.Flags().GetString("key")
		namespace, _ := cmd.Flags().GetString("namespace")
		hash, _ := cmd.Flags().GetString("hash")
		sigFile, _ := cmd.Flags().GetString("signature-file")

		pkc, err := ioutil.ReadFile(os.ExpandEnv(key))
		if err != nil {
			cmd.Println("FAIL: Unable to read your SSH private key. Make sure that you have entered its path correctly and have permission to access it.")
			return err
		}

		pk, err := ssh.ParsePrivateKey(pkc)
		if err != nil {
			cmd.Println("FAIL: Unable to parse your SSH private key. Make sure that it is a well-formatted SSH private key file.")
			return err
		}

		signer := sshsign.DefaultSigner(namespace, hash, pk)

		passes := true

		for _, file := range args {
			f, err := os.Open(file)
			if err != nil {
				cmd.Printf("FAIL: '%s' could not be opened for signing.\n", file)
				cmd.PrintErrln(err)
				passes = false
				continue
			}
			defer f.Close()

			sig, err := signer.Sign(f)
			if err != nil {
				cmd.Printf("FAIL: '%s' could not be signed.\n", file)
				cmd.PrintErrln(err)
				passes = false
				continue
			}

			armoured, err := sig.MarshalArmoured()
			if err != nil {
				cmd.Printf("FAIL: '%s' could not format the signature file correctly.\n", file)
				cmd.PrintErrln(err)
				passes = false
				continue
			}

			sigFile := strings.ReplaceAll(sigFile, "%f", file)
			if err := ioutil.WriteFile(sigFile, armoured, 0644); err != nil {
				cmd.Printf("FAIL: '%s' could not be saved.\n", sigFile)
				cmd.PrintErrln(err)
				passes = false
				continue
			}

			cmd.Printf("PASS: '%s' has been signed.\n", file)
		}

		if !passes {
			return fmt.Errorf("FAIL: One or more files could not be signed")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(signCmd)

	signCmd.Flags().StringP("key", "k", "$HOME/.ssh/id_rsa", "The SSH key to use for signing.")

	signCmd.Flags().StringP("namespace", "n", "file", "The namespace that the file should be signed for.")
	signCmd.Flags().String("signature-file", "%f.sig", "The signature file to use (defaults to the file.sig file).")
	signCmd.Flags().String("hash", "sha512", "The hash algorithm to use when verifying the file.")
}

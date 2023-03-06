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
	"os"

	"github.com/SierraSoftworks/shig/internal/core"
	"github.com/spf13/cobra"
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign",
	Short: "Sign a file using your SSH key.",
	Long: `Generate a cryptographic signature for a file using your SSH private key,
which can then be verified by anyone who knows your SSH public key. If you have your
key registered with GitHub, people can verify the signature using your GitHub username.`,
	Aliases: []string{"s"},
	Example: `shig sign file1.txt file2.txt`,
	Annotations: map[string]string{
		"Group": "Signing",
	},
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key, _ := cmd.Flags().GetString("key")
		namespace, _ := cmd.Flags().GetString("namespace")
		hash, _ := cmd.Flags().GetString("hash")
		sigFile, _ := cmd.Flags().GetString("signature-file")

		signer, err := core.NewSigner(cmd, key, namespace, hash, sigFile)
		if err != nil {
			cmd.PrintErrln(err)
			os.Exit(1)
		}

		passes := true

		for _, file := range args {
			if err := signer.Sign(file); err != nil {
				passes = false
				cmd.PrintErrln(err)
			}
		}

		if !passes {
			cmd.Println("FAIL: One or more files could not be signed")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(signCmd)

	signCmd.Flags().StringP("key", "k", "$HOME/.ssh/id_rsa", "The SSH key to use for signing.")

	signCmd.Flags().StringP("namespace", "n", "file", "The namespace that the file should be signed for.")
	signCmd.Flags().String("signature-file", "%f.sig", "The signature file to use (defaults to the file.sig file).")
	signCmd.Flags().String("hash", "sha512", "The hash algorithm to use when verifying the file.")
}

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
	"strings"

	"github.com/SierraSoftworks/shig/internal/core"
	"github.com/SierraSoftworks/shig/internal/publickeys"
	"github.com/spf13/cobra"
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify that the contents of a file match their signature.",
	Long: `This command allows you to verify that a file's contents have not been tampered
with and that they are signed by a trusted party.`,
	Aliases: []string{"v", "check"},
	Example: `shig verify --github notheotherben file1.txt file2.txt`,
	Args:    cobra.MinimumNArgs(1),
	Annotations: map[string]string{
		"Group": "Signing",
	},
	Run: func(cmd *cobra.Command, args []string) {
		namespace, _ := cmd.Flags().GetString("namespace")
		hash, _ := cmd.Flags().GetString("hash")
		sigFile, _ := cmd.Flags().GetString("signature-file")

		validator, err := getValidator(cmd)
		if err != nil {
			cmd.Println("FAIL: Unable to setup your requested validator. Make sure that you have entered the right value.")
			os.Exit(1)
		}

		if validator == nil {
			cmd.PrintErrln("WARNING: You have not provided a means of verifying the author of this signature. This makes it impossible to determine whether the file has been maliciously tampered with.")
		}

		verifier := core.NewVerifier(cmd, namespace, hash, sigFile, validator)

		passes := true

		for _, file := range args {
			if err := verifier.Verify(file); err != nil {
				passes = false
				cmd.PrintErrln(err)
			}
		}

		if !passes {
			cmd.Println("FAIL: One or more files failed verification")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().StringP("namespace", "n", "file", "The namespace that the file should be signed for.")
	verifyCmd.Flags().String("signature-file", "%f.sig", "The signature file to use (defaults to the file.sig file).")

	verifyCmd.Flags().String("hash", "sha512", "The hash algorithm to use when verifying the file.")
	verifyCmd.Flags().StringP("github", "G", "", "The GitHub username of the person you trust to sign the file.")
	verifyCmd.Flags().StringP("thumbprint", "T", "", "The fingerprint of the SSH key you trust to sign the file.")
	verifyCmd.Flags().StringP("publickey", "K", "", "The SSH public-key you trust to sign the file.")
}

func getValidator(cmd *cobra.Command) (publickeys.Validator, error) {
	github, _ := cmd.Flags().GetString("github")
	thumbprint, _ := cmd.Flags().GetString("thumbprint")
	ssh, _ := cmd.Flags().GetString("publickey")

	if github != "" {
		return publickeys.NewGitHubValidator(github)
	} else if thumbprint != "" {
		return publickeys.NewThumbprintValidator(thumbprint), nil
	} else if ssh != "" {
		if strings.HasPrefix(ssh, "ssh-") {
			return publickeys.NewSshKeyValidator(ssh)
		}

		return publickeys.NewAuthorizedKeysFileValidator(ssh)
	}

	return nil, nil
}

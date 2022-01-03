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

	"github.com/spf13/cobra"
)

var version = "dev"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "ssign",
	Short:   "Sign and verify files using your SSH key",
	Long:    `Use your existing SSH key(s) to sign and verify files, as a modern, easy to use, alternative to PGP.`,
	Version: version,
	Example: `# Signing files
shig sign file1.txt file2.txt
shig sign --key ~/.ssh/id_rsa file1.txt file2.txt

# Verifying files
shig verify file1.txt file2.txt --github notheotherben
shig verify file1.txt file2.txt --thumbprint 'SHA256:MW8+PD+j0wSkK8tY0hlk8868Ebl6jbmkwWPpgvhxEuk'
shig verify file1.txt file2.txt --publickey ~/.ssh/id_rsa.pub
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

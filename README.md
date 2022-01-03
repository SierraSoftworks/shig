# shig
**Cryptographically sign and verify files using SSH keys**

SSH keys are one of the most ubiquitous forms of public/private key cryptography
held by general computer users. Unlike PGP (which has several glaring issues),
X.509 (which is complex and costly to maintain), and bespoke systems like
minisign and signify, SSH keys are easy to use, familiar, and widely deployed.

This makes them exceptionally useful as a means of cryptographically signing and
verifying the content of files published in the open, and doubly so when used to
verify authorship against a user's GitHub profile keys.

While it is possible to do all of this using `ssh-keygen`, the flags required to
do so are far from intuitive. This command line application is designed to provide
a straightforward, easy-to-use, interface for signing and verifying files using
SSH keys in a manner that is fully compatible with `ssh-keygen`'s SSHSIG protocol.

## Features
 - Extremely quick and easy to use
 - Compatible with `ssh-keygen`'s SSHSIG protocol
 - Verify signed files using a user's GitHub profile keys

## Example

```bash
# Install the shig binary on your path (if you have Go installed)
> go install github.com/SierraSoftworks/shig

# Generate signatures for the provided files and write them to disk
> shig sign myfile.txt myotherfile.txt
PASS: 'myfile.txt' has been signed.
PASS: 'myotherfile.txt' has been signed.

# Verify the signatures for the provided files and trust @notheotherben's GitHub keys
> shig verify myfile.txt myotherfile.txt --github notheotherben
PASS: 'myfile.txt' has been signed by 'SHA256:MW8+PD+j0wSkK8tY0hlk8868Ebl6jbmkwWPpgvhxEuk'
PASS: 'myotherfile.txt' has been signed by 'SHA256:MW8+PD+j0wSkK8tY0hlk8868Ebl6jbmkwWPpgvhxEuk'
```
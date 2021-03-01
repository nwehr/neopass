package pgp

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/ssh/terminal"
)

type PGPDecrypter struct {
	SecretKeyring openpgp.EntityList
	GetPassphrase func() ([]byte, error)
}

func (d PGPDecrypter) Decrypt(text string) (string, error) {
	prompt := func(keys []openpgp.Key, symmetric bool) ([]byte, error) {
		for _, key := range keys {
			for identity := range key.Entity.Identities {
				fmt.Println(identity)
				break
			}

			passphrase, err := d.GetPassphrase()
			if err != nil {
				return nil, err
			}

			key.Entity.PrivateKey.Decrypt(passphrase)
			for _, subkey := range key.Entity.Subkeys {
				subkey.PrivateKey.Decrypt(passphrase)
			}
		}

		return nil, nil
	}

	password, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	md, err := openpgp.ReadMessage(bytes.NewBuffer(password), d.SecretKeyring, prompt, nil)
	if err != nil {
		return "", err
	}

	contents, err := ioutil.ReadAll(md.UnverifiedBody)
	return string(contents), err
}

func DefaultDecrypter(config Config) (PGPDecrypter, error) {
	keyringFile, err := os.Open(config.SecretKeyringPath)
	if err != nil {
		return PGPDecrypter{}, err
	}

	keyring, err := openpgp.ReadKeyRing(keyringFile)
	// keyring, err := openpgp.ReadArmoredKeyRing(keyringFile)
	if err != nil {
		return PGPDecrypter{}, err
	}

	return PGPDecrypter{
		SecretKeyring: keyring,
		GetPassphrase: func() ([]byte, error) {
			fmt.Print("passphrase: ")

			tty, err := os.Open("/dev/tty")
			if err != nil {
				return nil, err
			}

			defer tty.Close()
			defer fmt.Println()

			return terminal.ReadPassword(int(tty.Fd()))
		},
	}, err
}

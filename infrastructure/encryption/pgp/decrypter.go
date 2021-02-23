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
	Identity      string
	SecretKeyring openpgp.EntityList
	GetPassphrase func() ([]byte, error)
}

func (d PGPDecrypter) Decrypt(text string) (string, error) {
	entityByNameOrEmail := func(nameOrEmail string) (*openpgp.Entity, error) {
		for _, entity := range d.SecretKeyring {
			for _, identity := range entity.Identities {
				if nameOrEmail == identity.UserId.Name || nameOrEmail == identity.UserId.Email {
					return entity, nil
				}
			}
		}

		return nil, fmt.Errorf("entity %s does not exist", nameOrEmail)
	}

	entity, err := entityByNameOrEmail(d.Identity)
	if err != nil {
		return "", err
	}

	passphrase, err := d.GetPassphrase()
	if err != nil {
		return "", err
	}

	entity.PrivateKey.Decrypt(passphrase)
	for _, subkey := range entity.Subkeys {
		subkey.PrivateKey.Decrypt(passphrase)
	}

	password, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	md, err := openpgp.ReadMessage(bytes.NewBuffer(password), d.SecretKeyring, nil, nil)
	if err != nil {
		return "", err
	}

	contents, err := ioutil.ReadAll(md.UnverifiedBody)
	return string(contents), err
}

func DefaultDecrypter() (PGPDecrypter, error) {
	config, err := DefaultConfig()
	if err != nil {
		return PGPDecrypter{}, err
	}

	keyringFile, err := os.Open(config.SecretKeyringPath)
	if err != nil {
		return PGPDecrypter{}, err
	}

	keyring, err := openpgp.ReadKeyRing(keyringFile)

	return PGPDecrypter{
		Identity:      config.Identity,
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

package pgp

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"

	"golang.org/x/crypto/openpgp"
)

type PGPEncrypter struct {
	Identity      string
	PublicKeyring openpgp.EntityList
}

func (e PGPEncrypter) Encrypt(password string) (string, error) {
	entityByNameOrEmail := func(nameOrEmail string) (*openpgp.Entity, error) {
		for _, entity := range e.PublicKeyring {
			for _, identity := range entity.Identities {
				if nameOrEmail == identity.UserId.Name || nameOrEmail == identity.UserId.Email {
					return entity, nil
				}
			}
		}

		return nil, fmt.Errorf("entity %s does not exist", nameOrEmail)
	}

	entity, err := entityByNameOrEmail(e.Identity)
	if err != nil {
		return "", err
	}

	ciphertext := new(bytes.Buffer)
	plaintext, err := openpgp.Encrypt(ciphertext, openpgp.EntityList{entity}, nil, nil, nil)
	if err != nil {
		return "", err
	}

	_, err = plaintext.Write([]byte(password))
	if err != nil {
		return "", err
	}

	plaintext.Close()

	return base64.StdEncoding.EncodeToString(ciphertext.Bytes()), nil
}

func DefaultEncrypter(config Config) (PGPEncrypter, error) {
	keyringFile, err := os.Open(config.PublicKeyringPath)
	if err != nil {
		return PGPEncrypter{}, err
	}

	keyring, err := openpgp.ReadKeyRing(keyringFile)
	return PGPEncrypter{Identity: config.Identity, PublicKeyring: keyring}, err
}

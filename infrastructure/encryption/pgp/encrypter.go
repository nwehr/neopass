package pgp

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"

	"golang.org/x/crypto/openpgp"
)

type PGPEncrypter struct {
	Identities    []string
	PublicKeyring openpgp.EntityList
}

func (e PGPEncrypter) Encrypt(password string) (string, error) {
	keys := openpgp.EntityList{}

	for _, key := range e.PublicKeyring {
		for _, keyIdentity := range key.Identities {
			for _, identity := range e.Identities {
				if identity == keyIdentity.UserId.Name || identity == keyIdentity.UserId.Email {
					fmt.Printf("%s <%s>\n", keyIdentity.UserId.Name, keyIdentity.UserId.Email)
					keys = append(keys, key)
				}
			}
		}
	}

	ciphertext := new(bytes.Buffer)
	plaintext, err := openpgp.Encrypt(ciphertext, keys, nil, nil, nil)
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
	// keyring, err := openpgp.ReadArmoredKeyRing(keyringFile)
	return PGPEncrypter{Identities: config.Identities, PublicKeyring: keyring}, err
}

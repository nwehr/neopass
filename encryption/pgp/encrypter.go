package pgp

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/user"

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
		log.Fatal(err)
	}

	_, err = plaintext.Write([]byte(password))
	if err != nil {
		log.Fatal("plaintext.Write([]byte(password)) ", err)
	}

	plaintext.Close()

	return base64.StdEncoding.EncodeToString(ciphertext.Bytes()), nil
}

func DefaultEncrypter(identity string) (PGPEncrypter, error) {
	usr, _ := user.Current()

	keyringFile, err := os.Open(usr.HomeDir + "/.gnupg/pubring.gpg")
	if err != nil {
		return PGPEncrypter{}, err
	}

	keyring, err := openpgp.ReadKeyRing(keyringFile)
	return PGPEncrypter{Identity: identity, PublicKeyring: keyring}, err
}

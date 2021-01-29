package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/user"

	"golang.org/x/crypto/openpgp"
)

func main() {

	entity, err := entityByNameOrEmail(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	passphrase := []byte(os.Args[2])

	if entity.PrivateKey.Encrypted {
		if entity.PrivateKey != nil && entity.PrivateKey.Encrypted {
			err := entity.PrivateKey.Decrypt(passphrase)
			if err != nil {
				log.Fatalf("failed to decrypt key")
			}
		}
		for _, subkey := range entity.Subkeys {
			if subkey.PrivateKey != nil && subkey.PrivateKey.Encrypted {
				err := subkey.PrivateKey.Decrypt(passphrase)
				if err != nil {
					log.Fatalf("failed to decrypt subkey")
				}
			}
		}
	}

	ciphertext := new(bytes.Buffer)

	plaintext, err := openpgp.Encrypt(ciphertext, openpgp.EntityList{entity}, nil, nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	plaintext.Write([]byte(os.Args[3]))
	defer plaintext.Close()

	fmt.Println(base64.StdEncoding.EncodeToString(ciphertext.Bytes()))

}

func entityByNameOrEmail(nameOrEmail string) (*openpgp.Entity, error) {
	usr, _ := user.Current()

	keyringFile, err := os.Open(usr.HomeDir + "/.gnupg/secring.gpg")
	if err != nil {
		log.Fatal("os.Open(keyringPath): ", err)
	}

	entityList, err := openpgp.ReadKeyRing(keyringFile)
	if err != nil {
		log.Fatal("openpgp.ReadKeyRing(keyringFile): ", err)
	}

	for _, entity := range entityList {
		for _, identity := range entity.Identities {
			if nameOrEmail == identity.UserId.Name || nameOrEmail == identity.UserId.Email {
				return entity, nil
			}
		}
	}

	return nil, fmt.Errorf("entity does not exist")
}

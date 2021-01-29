package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/atotto/clipboard"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/ssh/terminal"
)

type Store struct {
	Identity string  `json:"identity"`
	Entries  Entries `json:"entries"`
}

func (s *Store) Add(name string, password string, keyring openpgp.EntityList) {
	entity, err := entityByNameOrEmail(s.Identity, keyring)
	if err != nil {
		log.Fatal("entityByNameOrEmail(store.Identity, kring) ", err)
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

	entry := Entry{
		Name:     name,
		Password: base64.StdEncoding.EncodeToString(ciphertext.Bytes()),
	}

	s.Entries = append(s.Entries, entry)
}

func (s Store) DecryptPassword(entry Entry, keyring openpgp.EntityList) {
	entity, err := entityByNameOrEmail(s.Identity, keyring)
	if err != nil {
		log.Fatal(`entityByNameOrEmail(store.Identity, kring) `, err)
	}

	fmt.Print("passphrase: ")
	passphrase, err := terminal.ReadPassword(int(os.Stdin.Fd()))

	entity.PrivateKey.Decrypt(passphrase)
	for _, subkey := range entity.Subkeys {
		subkey.PrivateKey.Decrypt(passphrase)
	}

	password, err := base64.StdEncoding.DecodeString(entry.Password)
	if err != nil {
		log.Fatal("base64.StdEncoding.DecodeString(entry.Password) ", err)
	}

	md, err := openpgp.ReadMessage(bytes.NewBuffer(password), keyring, nil, nil)
	if err != nil {
		log.Fatal("openpgp.ReadMessage(bytes.NewReader(password), kring, prompt, nil) ", err)
	}

	contents, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		log.Fatal("ioutil.ReadAll(md.UnverifiedBody) ", err)
	}

	if err = clipboard.WriteAll(string(contents)); err != nil {
		log.Fatal("clipboard.WriteAll(string(contents)) ", err)
	}

	fmt.Printf("\ncopied to clipboard\n")
}

type Entries []Entry

type Entry struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func entityByNameOrEmail(nameOrEmail string, kring openpgp.EntityList) (*openpgp.Entity, error) {
	for _, entity := range kring {
		for _, identity := range entity.Identities {
			if nameOrEmail == identity.UserId.Name || nameOrEmail == identity.UserId.Email {
				return entity, nil
			}
		}
	}

	return nil, fmt.Errorf("entity %s does not exist", nameOrEmail)
}

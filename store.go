package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
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
		panic(fmt.Sprintf("entityByNameOrEmail(store.Identity, kring) %s", err))
	}

	ciphertext := new(bytes.Buffer)
	plaintext, err := openpgp.Encrypt(ciphertext, openpgp.EntityList{entity}, nil, nil, nil)
	if err != nil {
		panic(fmt.Sprintf("openpgp.Encrypt(ciphertext, openpgp.EntityList{entity}, nil, nil, nil) %s", err))
	}

	_, err = plaintext.Write([]byte(password))
	if err != nil {
		panic(fmt.Sprintf("plaintext.Write([]byte(password)) %s", err))
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
		panic(fmt.Sprintf(`entityByNameOrEmail(store.Identity, kring) %s`, err))
	}

	fmt.Print("passphrase: ")

	tty, err := os.Open("/dev/tty")
	if err != nil {
		panic(fmt.Sprintf("could not open /dev/tty %s", err))
	}
	defer tty.Close()

	passphrase, err := terminal.ReadPassword(int(tty.Fd()))
	if err != nil {
		panic(fmt.Sprintf("terminal.ReadPassword(int(os.Stdin.Fd())) %s", err))
	}

	entity.PrivateKey.Decrypt(passphrase)
	for _, subkey := range entity.Subkeys {
		subkey.PrivateKey.Decrypt(passphrase)
	}

	password, err := base64.StdEncoding.DecodeString(entry.Password)
	if err != nil {
		panic(fmt.Sprintf("base64.StdEncoding.DecodeString(entry.Password) %s", err))
	}

	md, err := openpgp.ReadMessage(bytes.NewBuffer(password), keyring, nil, nil)
	if err != nil {
		panic(fmt.Sprintf("openpgp.ReadMessage(bytes.NewReader(password), kring, prompt, nil) %s", err))
	}

	contents, err := ioutil.ReadAll(md.UnverifiedBody)
	if err != nil {
		panic(fmt.Sprintf("ioutil.ReadAll(md.UnverifiedBody) %s", err))
	}

	if err = clipboard.WriteAll(string(contents)); err != nil {
		panic(fmt.Sprintf("clipboard.WriteAll(string(contents)) %s", err))
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

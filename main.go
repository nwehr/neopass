package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/atotto/clipboard"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	if len(os.Args) == 1 {
		listEntries()
	} else {
		switch os.Args[1] {
		case "init":
			createStore(os.Args[2])
		case "add":
			addEntry(os.Args[2], os.Args[3])
		default:
			findEntry(os.Args[1])
		}
	}
}

func createStore(nameOrEmail string) {
	usr, _ := user.Current()

	storeFile, err := os.Create(usr.HomeDir + "/.paws.json")
	if err != nil {
		log.Fatal("os.Create(storePath)", err)
	}

	store := Store{
		Identity: nameOrEmail,
	}

	err = json.NewEncoder(storeFile).Encode(store)
	if err != nil {
		log.Fatal("json.NewEncoder(storeFile).Encode(store)", err)
	}
}

func readStore() Store {
	usr, _ := user.Current()

	storeFile, err := os.Open(usr.HomeDir + "/.paws.json")
	if err != nil {
		log.Fatal("os.Open(storePath)", err)
	}

	store := Store{}

	err = json.NewDecoder(storeFile).Decode(&store)
	if err != nil {
		log.Fatal("json.NewDecoder(storeFile).Decode(&store)", err)
	}

	storeFile.Close()

	return store
}

func writeStore(store Store) {
	usr, _ := user.Current()

	storeFile, err := os.OpenFile(usr.HomeDir+"/.paws.json", os.O_RDWR, 0)
	if err != nil {
		log.Fatal("os.Open(storePath)", err)
	}

	defer storeFile.Close()

	storeFile.Truncate(0)
	storeFile.Seek(0, 0)

	encoded, err := json.MarshalIndent(store, "", "    ")

	_, err = storeFile.Write(encoded)
	if err != nil {
		log.Fatal(err)
	}
}

func listEntries() {
	store := readStore()

	for _, entry := range store.Entries {
		fmt.Println(entry.Name)
	}
}

func findEntry(name string) {
	store := readStore()

	for _, entry := range store.Entries {
		if name == entry.Name {
			kring, _ := secretKeyring()
			entity, err := entityByNameOrEmail(store.Identity, kring)
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

			md, err := openpgp.ReadMessage(bytes.NewBuffer(password), kring, nil, nil)
			if err != nil {
				log.Fatal("openpgp.ReadMessage(bytes.NewReader(password), kring, prompt, nil) ", err)
			}

			contents, err := ioutil.ReadAll(md.UnverifiedBody)
			if err != nil {
				log.Fatal("ioutil.ReadAll(md.UnverifiedBody) ", err)
			}

			clipboard.WriteAll(string(contents))
			fmt.Println("\npassword copped to clipboard")
		}
	}
}

func addEntry(name string, password string) {
	store := readStore()
	kring, _ := publicKeyring()
	entity, err := entityByNameOrEmail(store.Identity, kring)
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

	store.Entries = append(store.Entries, entry)

	writeStore(store)
}

func publicKeyring() (openpgp.EntityList, error) {
	usr, _ := user.Current()

	keyringFile, err := os.Open(usr.HomeDir + "/.gnupg/pubring.gpg")
	if err != nil {
		log.Fatal("os.Open(keyringPath): ", err)
	}

	return openpgp.ReadKeyRing(keyringFile)
}

func secretKeyring() (openpgp.EntityList, error) {
	usr, _ := user.Current()

	keyringFile, err := os.Open(usr.HomeDir + "/.gnupg/secring.gpg")
	if err != nil {
		log.Fatal("os.Open(keyringPath): ", err)
	}

	return openpgp.ReadKeyRing(keyringFile)
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

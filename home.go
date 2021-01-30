package main

import (
	"encoding/json"
	"log"
	"os"
	"os/user"

	"golang.org/x/crypto/openpgp"
)

type Home struct {
	HomeDir string
}

func NewHome() Home {
	usr, _ := user.Current()
	return Home{HomeDir: usr.HomeDir}
}

func (h Home) PublicKeyring() (openpgp.EntityList, error) {
	keyringFile, err := os.Open(h.HomeDir + "/.gnupg/pubring.gpg")
	if err != nil {
		log.Fatal("os.Open(keyringPath): ", err)
	}

	return openpgp.ReadKeyRing(keyringFile)
}

func (h Home) SecretKeyring() (openpgp.EntityList, error) {
	keyringFile, err := os.Open(h.HomeDir + "/.gnupg/secring.gpg")
	if err != nil {
		log.Fatal("os.Open(keyringPath): ", err)
	}

	return openpgp.ReadKeyRing(keyringFile)
}

func (h Home) ReadStore() Store {
	storeFile, err := os.Open(h.HomeDir + "/.paws/store.json")
	if err != nil {
		log.Fatal("os.Open(storePath)", err)
	}

	defer storeFile.Close()

	store := Store{}

	err = json.NewDecoder(storeFile).Decode(&store)
	if err != nil {
		log.Fatal("json.NewDecoder(storeFile).Decode(&store)", err)
	}

	return store
}

func (h Home) WriteStore(store Store) {
	if _, err := os.Stat(h.HomeDir + "/.paws"); os.IsNotExist(err) {
		os.Mkdir(h.HomeDir+"/.paws", os.ModeDir|os.ModePerm)
	}

	storeFile, err := os.OpenFile(h.HomeDir+"/.paws/store.json", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatal("os.Open(storePath)", err)
	}

	defer storeFile.Close()

	encoded, err := json.MarshalIndent(store, "", "    ")
	if err != nil {
		log.Fatal(`json.MarshalIndent(store, "", "    ") `, err)
	}

	storeFile.Truncate(0)
	storeFile.Seek(0, 0)

	if _, err = storeFile.Write(encoded); err != nil {
		log.Fatal(err)
	}
}

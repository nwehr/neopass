package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nwehr/paws/application/commands"
	"github.com/nwehr/paws/application/queries"
	"github.com/nwehr/paws/encryption/pgp"
	"github.com/nwehr/paws/persistance/file"
	"golang.org/x/crypto/ssh/terminal"
)

func getPassword() ([]byte, error) {
	tty, err := os.Open("/dev/tty")
	if err != nil {
		return nil, err
	}

	defer tty.Close()

	return terminal.ReadPassword(int(tty.Fd()))
}

func main() {
	r := file.DefaultStoreRepository()

	if len(os.Args) == 1 {
		entries, err := queries.ListEntries{}.Execute(r)
		if err != nil {
			log.Fatal(err)
		}

		for _, entryName := range entries {
			fmt.Println(entryName)
		}
		return
	}

	switch os.Args[1] {
	case "add":
		store, _ := r.Load()
		enc, err := pgp.DefaultEncrypter(store.Identity)
		if err != nil {
			log.Fatal(err)
		}

		name := os.Args[2]
		password, _ := getPassword()

		err = commands.AddEntry{Name: name, Password: string(password)}.Execute(enc, r)
		if err != nil {
			log.Fatal(err)
		}
	default:
		store, _ := r.Load()
		d, err := pgp.DefaultDecrypter(store.Identity)
		if err != nil {
			log.Fatal(err)
		}

		name := os.Args[1]
		password, err := queries.GetEntry{Name: name}.Execute(d, r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println()
		fmt.Println(password)
	}
}

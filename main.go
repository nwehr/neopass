package main

import (
	"fmt"
	"os"

	"github.com/nwehr/paws/application/commands"
	"github.com/nwehr/paws/application/queries"
	"github.com/nwehr/paws/encryption/pgp"
	"github.com/nwehr/paws/persistance"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	p := persistance.DefaultFilePersister()

	if len(os.Args) == 1 {
		names, err := queries.AllEntryNames{}.Execute(p)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, name := range names {
			fmt.Println(name)
		}
		return
	}

	switch os.Args[1] {
	case "add":
		store, _ := p.Load()
		enc, err := pgp.DefaultEncrypter(store.Identity)
		if err != nil {
			fmt.Println(err)
			return
		}

		name := os.Args[2]
		password, _ := getPassword()

		err = commands.AddEntry{Name: name, Password: string(password)}.Execute(enc, p)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		store, _ := p.Load()
		d, err := pgp.DefaultDecrypter(store.Identity)
		if err != nil {
			fmt.Println(err)
			return
		}

		name := os.Args[1]
		password, err := queries.GetEntryPassword{Name: name}.Execute(d, p)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(password)
	}
}

func getPassword() ([]byte, error) {
	fmt.Print("password: ")

	tty, err := os.Open("/dev/tty")
	if err != nil {
		return nil, err
	}

	defer tty.Close()
	defer fmt.Println()

	return terminal.ReadPassword(int(tty.Fd()))
}

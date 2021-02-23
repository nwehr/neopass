package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nwehr/paws/application/commands"
	"github.com/nwehr/paws/application/queries"
	"github.com/nwehr/paws/encryption/pgp"
	"github.com/nwehr/paws/persistance"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	p := persistance.DefaultFilePersister()
	store, _ := p.Load()
	enc, err := pgp.DefaultEncrypter(store.Identity)
	dec, err := pgp.DefaultDecrypter(store.Identity)

	if weAreInAPipe() {
		name, err := readNameFromPipe()
		if err != nil {
			fmt.Println(err)
			return
		}
		password, err := queries.GetEntryPassword{Name: name}.Execute(dec, p)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(password)
		return
	}

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
		name := os.Args[2]
		password, _ := getPassword()

		err = commands.AddEntry{Name: name, Password: string(password)}.Execute(enc, p)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "rm":
		name := os.Args[2]
		err = commands.RemoveEntry{name}.Execute(p)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		name := os.Args[1]
		password, err := queries.GetEntryPassword{Name: name}.Execute(dec, p)
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

func weAreInAPipe() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	if info.Mode()&os.ModeCharDevice == 0 {
		return true
	}

	return false
}

func readNameFromPipe() (string, error) {
	r := bufio.NewReader(os.Stdin)

	var output []rune

	for {
		input, _, err := r.ReadRune()
		if err != nil && err == io.EOF {
			r.ReadRune()
			break
		}
		output = append(output, input)
	}

	runesToString := func(runes []rune) (outString string) {
		for _, v := range runes {
			outString += string(v)
		}
		return
	}

	name := strings.TrimSpace(runesToString(output))
	return name, nil

}

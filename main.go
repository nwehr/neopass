package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"
	"strings"

	"github.com/nwehr/paws/core/domain"
	"github.com/nwehr/paws/core/usecases"
	"github.com/nwehr/paws/infrastructure/encryption/pgp"
	"github.com/nwehr/paws/infrastructure/persistance"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	p := persistance.DefaultFilePersister()

	if weAreInAPipe() {
		dec, err := pgp.DefaultDecrypter()
		if err != nil {
			fmt.Println(err)
			return
		}

		name, err := nameFromPipe()
		if err != nil {
			fmt.Println(err)
			return
		}
		password, err := usecases.GetDecryptedPassword{p, dec}.Run(name)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(password)
		return
	}

	if len(os.Args) == 1 {
		names, err := usecases.GetAllEntryNames{p}.Run()
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
	case "init":
		identity := os.Args[2]

		usr, _ := user.Current()

		config := pgp.Config{
			Identity:          identity,
			PublicKeyringPath: usr.HomeDir + "/.gnupg/pubring.gpg",
			SecretKeyringPath: usr.HomeDir + "/.gnupg/secring.gpg",
		}

		if err := pgp.SaveConfig(config); err != nil {
			fmt.Println(err)
			return
		}

		_, err := p.Load()
		if err != nil {
			p.Save(domain.Store{})
		}
	case "add":
		enc, err := pgp.DefaultEncrypter()
		if err != nil {
			fmt.Println(err)
			return
		}

		name := os.Args[2]
		password, _ := ttyPassword()

		err = usecases.AddEntry{enc, p}.Run(name, string(password))
		if err != nil {
			fmt.Println(err)
			return
		}
	case "rm":
		name := os.Args[2]
		err := usecases.RemoveEntry{p}.Run(name)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		dec, err := pgp.DefaultDecrypter()
		if err != nil {
			fmt.Println(err)
			return
		}

		name := os.Args[1]
		password, err := usecases.GetDecryptedPassword{p, dec}.Run(name)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(password)
	}
}

func ttyPassword() ([]byte, error) {
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

func nameFromPipe() (string, error) {
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

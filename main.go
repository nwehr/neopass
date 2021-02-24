package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	info, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println("os.Stdin.Stat() ", err)
		os.Exit(1)
	}

	if info.Mode()&os.ModeCharDevice == 0 {
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

		home := NewHome()

		keyring, _ := home.SecretKeyring()
		store := home.ReadStore()

		for _, entry := range store.Entries {
			if name == entry.Name {
				store.DecryptPassword(entry, keyring)
			}
		}

		return
	}

	home := NewHome()

	if len(os.Args) == 1 {
		for _, entry := range home.ReadStore().Entries {
			fmt.Println(entry.Name)
		}

		return
	}

	switch os.Args[1] {
	case "init":
		store := Store{Identity: os.Args[2]}
		home.WriteStore(store)
	case "add":
		keyring, _ := home.PublicKeyring()
		store := home.ReadStore()

		fmt.Print("password: ")

		tty, err := os.Open("/dev/tty")
		if err != nil {
			panic(fmt.Sprintf("could not open /dev/tty %s", err))
		}
		defer tty.Close()

		password, err := terminal.ReadPassword(int(tty.Fd()))
		if err != nil {
			panic(fmt.Sprintf("terminal.ReadPassword(int(os.Stdin.Fd())) %s", err))
		}

		fmt.Println()

		store.Add(os.Args[2], string(password), keyring)
		home.WriteStore(store)
	case "rm":
		store := home.ReadStore()

		for i, entry := range store.Entries {
			if os.Args[2] == entry.Name {
				store.Entries = append(store.Entries[:i], store.Entries[i+1:]...)
			}
		}

		home.WriteStore(store)
	default:
		keyring, _ := home.SecretKeyring()
		store := home.ReadStore()

		for _, entry := range store.Entries {
			if os.Args[1] == entry.Name {
				store.DecryptPassword(entry, keyring)
			}
		}
	}
}

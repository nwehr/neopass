package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
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

		store.Add(os.Args[2], os.Args[3], keyring)
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

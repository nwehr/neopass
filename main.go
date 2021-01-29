package main

import (
	"fmt"
	"os"
)

func main() {
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

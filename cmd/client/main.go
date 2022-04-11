package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/nwehr/neopass"

	"github.com/nwehr/neopass/cmd/client/cmd/add"
	"github.com/nwehr/neopass/cmd/client/cmd/gen"
	"github.com/nwehr/neopass/cmd/client/cmd/initstore"
	"github.com/nwehr/neopass/cmd/client/cmd/ls"
	"github.com/nwehr/neopass/cmd/client/cmd/rm"
	"github.com/nwehr/neopass/pkg/config"
	enc "github.com/nwehr/neopass/pkg/encryption/age"
	"github.com/nwehr/neopass/pkg/repos"
)

var (
	Version string
	Built   string
)

func main() {
	if len(os.Args) == 1 {
		opts, err := config.GetConfigOptions(os.Args)
		if err != nil {
			Fatalf("%v\n", err)
		}

		err = ls.RunLs(opts)
		if err != nil {
			Fatalf("%v\n", err)
		}

		os.Exit(0)
	}

	switch os.Args[1] {
	case "init":
		opts, err := initstore.GetInitOptions(os.Args)
		if err != nil {
			Fatalf("could not get init options: %v\n", err)
		}

		err = initstore.RunInit(opts)
		if err != nil {
			Fatalf("%v\n", err)
		}

	case "-h":
		fallthrough
	case "--help":
		fallthrough
	case "help":
		fmt.Println("Usage")
		fmt.Println("  neopass [<command> <name>] | [<name>]")
		fmt.Println()
		fmt.Println("  Commands")
		fmt.Println("    init [--piv [slot]] [--neopass.cloud]  Initialize store")
		fmt.Println("    add   name           Add entry identified by name")
		fmt.Println("    gen   name           Generate new entry identified by name")
		fmt.Println("    rm    name           Remove entry identified by name")
		fmt.Println("    store [name]         Switch to store identified by name or list stores")
		fmt.Println()
		fmt.Println("  Examples")
		fmt.Println("     Initialize new password store on neopass cloud")
		fmt.Println("         neopass init --piv --neopass.cloud")
		fmt.Println()
		fmt.Println("     Add a new entry for github.com")
		fmt.Println("         neopass add github.com")
		fmt.Println()
		fmt.Println("     Switch to a password store named default")
		fmt.Println("         neopass store default")

	case "version":
		fmt.Printf("{\"version\":\"%s\", \"built\":\"%s\"}", Version, Built)

	case "import":
		opts, err := config.GetConfigOptions(os.Args)
		if err != nil {
			Fatalf("could not load config: %v\n", err)
		}

		store, err := opts.Config.GetCurrentStore()
		if err != nil {
			Fatalf("could not get current store: %v\n", err)
		}

		r := repos.FileRepo{Path: store.Location}

		file, err := os.Open(os.Args[2])
		if err != nil {
			Fatalf("could not open csv file: %v\n", err)
		}

		rows, err := csv.NewReader(file).ReadAll()
		if err != nil {
			Fatalf("could not parse csv file: %v\n", err)
		}

		for _, row := range rows {
			name := strings.TrimSpace(row[0])
			plain := strings.TrimSpace(row[1])

			enc, err := enc.NewAgeEncrypter(store.Age.Recipients)
			if err != nil {
				Fatalf("could not get encrypter: %v\n", err)
			}

			encrypted, err := enc.Encrypt(string(plain))
			if err != nil {
				Fatalf("could not encrypt password: %v\n", err)
			}

			entry := neopass.Entry{
				Name:     name,
				Password: encrypted,
			}

			if err := r.AddEntry(entry); err != nil {
				Fatalf("could not add entry: %v\n", err)
			}
		}

	case "add":
		opts, err := add.GetAddOptions(os.Args)
		if err != nil {
			Fatalf("%v\n", err)
		}

		err = add.RunAdd(opts)
		if err != nil {
			Fatalf("%v\n", err)
		}

	case "gen":
		opts, err := gen.GetGenOptions(os.Args)
		if err != nil {
			Fatalf("%v\n", err)
		}

		err = gen.RunGen(opts)
		if err != nil {
			Fatalf("%v\n", err)
		}

	case "rm":
		opts, err := rm.GetRmOptions(os.Args)
		if err != nil {
			Fatalf("%v\n", err)
		}

		err = rm.RunRm(opts)
		if err != nil {
			Fatalf("%v\n", err)
		}

	case "store":
		opts, err := config.GetConfigOptions(os.Args)
		if err != nil {
			Fatalf("could not load config: %v\n", err)
		}

		for _, store := range opts.Config.Stores {
			marker := "  "

			if store.Name == opts.Config.CurrentStore {
				marker = "->"
			}

			fmt.Printf("%s %s\n", marker, store.Name)
		}
	default:
		opts, err := config.GetConfigOptions(os.Args)
		if err != nil {
			Fatalf("could not load config: %v\n", err)
		}

		store, err := opts.Config.GetCurrentStore()
		if err != nil {
			Fatalf("could not get current store: %v\n", err)
		}

		r, err := opts.Config.GetCurrentRepo()

		entry, err := r.GetEntryByName(os.Args[1])
		if err != nil {
			Fatalf("could not find entry: %v\n", err)
		}

		identity, err := store.Age.UnlockIdentity()
		if err != nil {
			Fatalf("could not unlock identity: %v\n", err)
		}

		dec, err := enc.NewAgeDecrypter(identity)
		if err != nil {
			Fatalf("could not get decrypter: %v\n", err)
		}

		decrypted, err := dec.Decrypt(entry.Password)
		if err != nil {
			Fatalf("could not decrypt password: %v\n", err)
		}

		if err = clipboard.WriteAll(decrypted); err != nil {
			Fatalf("coult not write password to clipboard: %v", err)
		}

		fmt.Println("copied to clipboard")
	}
}

func Fatalf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

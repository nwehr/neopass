package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	mrand "math/rand"

	"github.com/atotto/clipboard"
	"github.com/nwehr/npass"
	"golang.org/x/crypto/ssh/terminal"

	enc "github.com/nwehr/npass/pkg/encryption/age"
	"github.com/nwehr/npass/pkg/repos"
)

func config() (npass.Config, error) {
	c := npass.Config{}
	err := c.ReadFile(npass.DefaultConfigFile)

	return c, err
}

func main() {
	if len(os.Args) == 1 {
		c, err := config()
		if err != nil {
			Fatalf("could not load config: %v\n", err)
		}

		store, err := c.GetCurrentStore()
		if err != nil {
			Fatalf("could not get current store: %v\n", err)
		}

		r := repos.FileRepo{Path: store.Location}

		names, err := r.ListEntryNames()
		if err != nil {
			Fatalf("could not list entry names: %v\n", err)
		}

		for _, name := range names {
			fmt.Println(name)
		}

		os.Exit(0)
	}

	switch os.Args[1] {
	case "-i":
		fallthrough
	case "--init":
		newAgeConfig := func() (npass.AgeConfig, error) {
			if len(os.Args) > 2 && os.Args[2] == "--piv" {
				var slotAddr uint32 = 0x9e

				if len(os.Args) > 3 {
					addr, err := strconv.ParseUint(os.Args[3], 16, 64)
					if err != nil {
						return npass.AgeConfig{}, err
					}

					slotAddr = uint32(addr)
				}

				return npass.NewPIVAgeConfig(slotAddr)
			}

			return npass.NewDefaultAgeConfig()
		}

		ageConfig, err := newAgeConfig()
		if err != nil {
			fmt.Printf("could not setup initial store: %v", err)
			os.Exit(1)
		}

		c := npass.Config{
			CurrentStore: "default",
			Stores: []npass.StoreConfig{
				{
					Name:     "default",
					Location: npass.DefaultStoreFile,
					Age:      ageConfig,
				},
			},
		}

		if err := c.WriteFile(npass.DefaultConfigFile); err != nil {
			fmt.Printf("could not write config to file: %v", err)
			os.Exit(1)
		}

	case "-h":
		fallthrough
	case "--help":
		fmt.Println("Usage")
		fmt.Println("  npass [<option> <name>] | [<name>]")
		fmt.Println()
		fmt.Println("  Options")
		fmt.Println("    -i | --init [--piv [slot]]  Setup initial store optionaly protected with a security card")
		fmt.Println("    -a | --add   name           Add entry identified by name")
		fmt.Println("    -g | --gen   name           Generate new entry identified by name")
		fmt.Println("    -r | --rm    name           Remove entry identified by name")
		fmt.Println("    -s | --store name           Switch to store identified by name or list stores")
		fmt.Println()
		fmt.Println("  Examples")
		fmt.Println("     Add a new entry for github.com")
		fmt.Println("         npass -a github.com")
		fmt.Println()
		fmt.Println("     Switch to a password store named default")
		fmt.Println("         npass -s default")

	case "-a":
		fallthrough
	case "--add":
		c, err := config()
		if err != nil {
			Fatalf("could not load config: %v\n", err)
		}

		store, err := c.GetCurrentStore()
		if err != nil {
			Fatalf("could not get current store: %v\n", err)
		}

		r := repos.FileRepo{Path: store.Location}

		plain, err := ttyPassword()
		if err != nil {
			Fatalf("could not get password: %v\n", err)
		}

		enc, err := enc.NewAgeEncrypter(store.Age.Recipients)
		if err != nil {
			Fatalf("could not get encrypter: %v\n", err)
		}

		encrypted, err := enc.Encrypt(string(plain))
		if err != nil {
			Fatalf("could not encrypt password: %v\n", err)
		}

		entry := npass.Entry{
			Name:     os.Args[2],
			Password: encrypted,
		}

		if err := r.AddEntry(entry); err != nil {
			Fatalf("could not add entry: %v\n", err)
		}

	case "-g":
		fallthrough
	case "--gen":
		c, err := config()
		if err != nil {
			Fatalf("could not load config: %v\n", err)
		}

		store, err := c.GetCurrentStore()
		if err != nil {
			Fatalf("could not get current store: %v\n", err)
		}

		r := repos.FileRepo{Path: store.Location}

		plain := genPassword()

		enc, err := enc.NewAgeEncrypter(store.Age.Recipients)
		if err != nil {
			Fatalf("could not get encrypter: %v\n", err)
		}

		encrypted, err := enc.Encrypt(string(plain))
		if err != nil {
			Fatalf("could not encrypt password: %v\n", err)
		}

		entry := npass.Entry{
			Name:     os.Args[2],
			Password: encrypted,
		}

		if err := r.AddEntry(entry); err != nil {
			Fatalf("could not add entry: %v\n", err)
		}

		if err = clipboard.WriteAll(plain); err != nil {
			Fatalf("coult not write password to clipboard: %v", err)
		}

		fmt.Println("copied to clipboard")

	case "-r":
		fallthrough
	case "--rm":
		c, err := config()
		if err != nil {
			Fatalf("could not load config: %v\n", err)
		}

		store, err := c.GetCurrentStore()
		if err != nil {
			Fatalf("coult not get current store: %v", err)
		}

		r := repos.FileRepo{Path: store.Location}

		if err = r.RemoveEntryByName(os.Args[2]); err != nil {
			Fatalf("could not remove entry : %v", err)
		}

	case "-s":
		c, err := config()
		if err != nil {
			Fatalf("could not load config: %v\n", err)
		}

		for _, store := range c.Stores {
			marker := "  "

			if store.Name == c.CurrentStore {
				marker = "->"
			}

			fmt.Printf("%s %s\n", marker, store.Name)
		}
	default:
		c, err := config()
		if err != nil {
			Fatalf("could not load config: %v\n", err)
		}

		store, err := c.GetCurrentStore()
		if err != nil {
			Fatalf("could not get current store: %v\n", err)
		}

		r := repos.FileRepo{Path: store.Location}

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

func genPassword() string {
	mrand.Seed(time.Now().UnixNano())

	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	segment := func() string {
		segment := []byte{}

		for {
			segment = append(segment, chars[mrand.Intn(len(chars))])

			if len(segment) == 4 {
				break
			}
		}

		return string(segment)
	}

	return fmt.Sprintf("%s+%s+%s+%s", segment(), segment(), segment(), segment())
}

func Fatalf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(1)
}

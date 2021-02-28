package cli

import (
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"time"

	"github.com/atotto/clipboard"
	"github.com/nwehr/npass/core/domain"
	"github.com/nwehr/npass/core/usecases"
	"github.com/nwehr/npass/infrastructure/encryption"
	"github.com/nwehr/npass/infrastructure/encryption/pgp"
	"golang.org/x/crypto/ssh/terminal"
)

type Cli struct {
	Repository domain.Repository
	Encrypter  encryption.Encrypter
	Decrypter  encryption.Decrypter
}

func (c Cli) Start(args []string) {
	if len(args) == 1 {
		if weAreInAPipe() {
			c.Get()
			return
		}

		c.List()
		return
	}

	switch args[1] {
	case "-i":
		fallthrough
	case "init":
		c.Init(args[2])

	case "-a":
		fallthrough
	case "add":
		c.Add()

	case "-g":
		fallthrough
	case "gen":
		c.Gen()

	case "-r":
		fallthrough
	case "rm":
		c.Rm()

	case "-c":
		fallthrough
	case "ctx":
		c.Ctx()

	default:
		c.Get()
	}
}

func (Cli) Init(identity string) {
	usr, _ := user.Current()

	conf := Config{
		CurrentContext: "default",
		Contexts: []Context{
			{
				Name: "default",
				Pgp: &pgp.Config{
					Identities:        []string{identity},
					PublicKeyringPath: usr.HomeDir + "/.gnupg/pubring.gpg",
					SecretKeyringPath: usr.HomeDir + "/.gnupg/secring.gpg",
				},
			},
		},
	}

	if err := SaveConfig(conf); err != nil {
		fmt.Println(err)
		return
	}
}

func (c Cli) Add() {
	name := os.Args[2]
	password, _ := ttyPassword()

	err := usecases.AddEntry{c.Repository, c.Encrypter}.Run(name, string(password))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (c Cli) Gen() {
	name := os.Args[2]
	password := generatePassword()

	err := usecases.AddEntry{c.Repository, c.Encrypter}.Run(name, string(password))
	if err != nil {
		fmt.Println(err)
		return
	}

	clipboard.WriteAll(password)
	fmt.Println("copied to clipboard")
}

func (c Cli) Rm() {
	name := os.Args[2]
	err := usecases.RemoveEntry{c.Repository}.Run(name)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (c Cli) Get() {
	name := ""

	if weAreInAPipe() {
		name, _ = readFromPipe()
	} else {
		name = os.Args[1]
	}

	password, err := usecases.GetDecryptedPassword{c.Repository, c.Decrypter}.Run(name)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = clipboard.WriteAll(password)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("copied to clipboard")
}

func (c Cli) List() {
	names, err := usecases.GetAllEntryNames{c.Repository}.Run()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, name := range names {
		fmt.Println(name)
	}
	return
}

func (c Cli) Ctx() {
	conf, err := LoadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(os.Args) == 3 {
		conf.CurrentContext = os.Args[2]
		SaveConfig(conf)
		return
	}

	for _, ctx := range conf.Contexts {
		marker := "  "

		if ctx.Name == conf.CurrentContext {
			marker = "->"
		}

		fmt.Printf("%s %s\n", marker, ctx.Name)
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

func generatePassword() string {
	rand.Seed(time.Now().UnixNano())

	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	segment := func() string {
		segment := []byte{}

		for {
			segment = append(segment, chars[rand.Intn(len(chars))])

			if len(segment) == 4 {
				break
			}
		}

		return string(segment)
	}

	return fmt.Sprintf("%s-%s-%s-%s", segment(), segment(), segment(), segment())
}

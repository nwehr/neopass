package cli

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"strings"
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
			name, _ := readFromPipe()
			password, err := c.Get(name)
			if err != nil {
				fmt.Println(err)
				return
			}

			if clipboard.Unsupported {
				fmt.Println(password)
				return
			}

			clipboard.WriteAll(password)
			fmt.Println("copied to clipboard")

			return
		}

		names, err := c.List()
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, name := range names {
			fmt.Println(name)
		}
		return
	}

	switch args[1] {
	case "-i":
		fallthrough
	case "--init":
		c.Init()

	case "-a":
		fallthrough
	case "--add":
		name := os.Args[2]
		password, _ := ttyPassword()

		if err := c.Add(name, string(password)); err != nil {
			fmt.Println(err)
			return
		}

	case "-g":
		fallthrough
	case "--gen":
		c.Gen()

	case "-r":
		fallthrough
	case "--rm":
		if err := c.Rm(os.Args[2]); err != nil {
			fmt.Println(err)
			return
		}

	case "-c":
		fallthrough
	case "--context":
		c.Ctx()

	case "--reencrypt":
		names, err := c.List()
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, name := range names {
			password, err := c.Get(name)
			if err != nil {
				fmt.Println(err)
				break
			}

			if err = c.Rm(name); err != nil {
				fmt.Println(err)
				break
			}

			if err = c.Add(name, password); err != nil {
				fmt.Println(err)
				break
			}
		}

	default:
		password, err := c.Get(os.Args[1])
		if err != nil {
			fmt.Println(err)
			return
		}

		if clipboard.Unsupported {
			fmt.Println(password)
			return
		}

		clipboard.WriteAll(password)
		fmt.Println("copied to clipboard")
	}
}

func (Cli) Init() {
	usr, _ := user.Current()

	contextName, _ := ttyPrompt("context", "default")
	identity, _ := ttyPrompt("identity", "")
	storeLocation, _ := ttyPrompt("store", usr.HomeDir+"/.npass/store.json")
	authToken, _ := ttyPrompt("auth", "")
	pubring, _ := ttyPrompt("pubring", usr.HomeDir+"/.gnupg/pubring.gpg")
	secring, _ := ttyPrompt("secring", usr.HomeDir+"/.gnupg/secring.gpg")

	conf, _ := LoadConfig()
	conf.CurrentContext = contextName
	conf.Contexts = append(conf.Contexts, Context{
		Name:          contextName,
		StoreLocation: &storeLocation,
		AuthToken:     &authToken,
		Pgp: &pgp.Config{
			Identities:        []string{identity},
			PublicKeyringPath: pubring,
			SecretKeyringPath: secring,
		},
	})

	if err := SaveConfig(conf); err != nil {
		fmt.Println("could not save config:", err)
		return
	}
}

func (c Cli) Add(name, password string) error {
	return usecases.AddEntry{c.Repository, c.Encrypter}.Run(name, string(password))
}

func (c Cli) Gen() {
	name := os.Args[2]
	password := generatePassword()

	err := usecases.AddEntry{c.Repository, c.Encrypter}.Run(name, string(password))
	if err != nil {
		fmt.Println(err)
		return
	}

	if clipboard.Unsupported {
		fmt.Println(password)
		return
	}

	clipboard.WriteAll(password)
	fmt.Println("copied to clipboard")
}

func (c Cli) Rm(name string) error {
	return usecases.RemoveEntry{c.Repository}.Run(name)
}

func (c Cli) Get(name string) (string, error) {
	return usecases.GetDecryptedPassword{c.Repository, c.Decrypter}.Run(name)
}

func (c Cli) List() ([]string, error) {
	return usecases.GetAllEntryNames{c.Repository}.Run()
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

func ttyPrompt(prompt, defaultValue string) (string, error) {
	if defaultValue != "" {
		prompt += " [" + defaultValue + "]"
	}

	fmt.Print(prompt + ": ")

	value, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if value == "\n" {
		value = defaultValue
	}

	return strings.TrimSpace(value), err
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

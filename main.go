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
	"github.com/nwehr/paws/infrastructure/config"
	"github.com/nwehr/paws/infrastructure/encryption"
	"github.com/nwehr/paws/infrastructure/encryption/pgp"
	"github.com/nwehr/paws/infrastructure/persistance"
	"github.com/nwehr/paws/interface/api"
	"github.com/nwehr/paws/interface/cli"

	"golang.org/x/crypto/ssh/terminal"
)

func repo() (domain.StoreRepository, error) {
	conf, err := config.LoadCliConfig()
	var repo domain.StoreRepository

	if conf.StoreLocation != nil && strings.HasPrefix(*conf.StoreLocation, "http") {
		repo = persistance.ApiRepository{*conf.StoreLocation, *conf.AuthToken}
	} else {
		repo = persistance.DefaultFileRepository()
	}

	return repo, err
}

func main() {

	if weAreInAPipe() {
		conf, err := config.LoadPgpConfig()
		if err != nil {
			fmt.Println(err)
			return
		}

		dec, err := pgp.DefaultDecrypter(conf)
		if err != nil {
			fmt.Println(err)
			return
		}

		name, err := nameFromPipe()
		if err != nil {
			fmt.Println(err)
			return
		}

		repo, err := repo()
		if err != nil {
			fmt.Println(err)
			return
		}

		password, err := usecases.GetDecryptedPassword{repo, dec}.Run(name)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(password)
		return
	}

	if len(os.Args) == 1 {
		repo, err := repo()
		if err != nil {
			fmt.Println(err)
			return
		}

		names, err := usecases.GetAllEntryNames{repo}.Run()
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

		conf := pgp.Config{
			Identity:          identity,
			PublicKeyringPath: usr.HomeDir + "/.gnupg/pubring.gpg",
			SecretKeyringPath: usr.HomeDir + "/.gnupg/secring.gpg",
		}

		if err := config.SavePgpConfig(conf); err != nil {
			fmt.Println(err)
			return
		}

		if err := config.SaveCliConfig(cli.Config{}); err != nil {
			fmt.Println(err)
			return
		}

		if err := config.SaveApiConfig(api.Config{}); err != nil {
			fmt.Println(err)
			return
		}
	case "add":
		pgpConfig, err := config.LoadPgpConfig()
		if err != nil {
			fmt.Println(err)
			return
		}

		repo, err := repo()
		if err != nil {
			fmt.Println(err)
			return
		}

		enc, err := pgp.DefaultEncrypter(pgpConfig)
		if err != nil {
			fmt.Println(err)
			return
		}

		name := os.Args[2]
		password, _ := ttyPassword()

		err = usecases.AddEntry{repo, enc}.Run(name, string(password))
		if err != nil {
			fmt.Println(err)
			return
		}
	case "rm":
		repo, err := repo()
		if err != nil {
			fmt.Println(err)
			return
		}

		name := os.Args[2]
		err = usecases.RemoveEntry{repo}.Run(name)
		if err != nil {
			fmt.Println(err)
			return
		}
	case "start-server":
		apiConfig, err := config.LoadApiConfig()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(api.Api{apiConfig, encryption.NoEncrypter{}, encryption.NoDecrypter{}}.Start())
		return
	default:
		repo, err := repo()
		if err != nil {
			fmt.Println(err)
			return
		}

		pgpConfig, err := config.LoadPgpConfig()
		if err != nil {
			fmt.Println(err)
			return
		}

		dec, err := pgp.DefaultDecrypter(pgpConfig)
		if err != nil {
			fmt.Println(err)
			return
		}

		name := os.Args[1]
		password, err := usecases.GetDecryptedPassword{repo, dec}.Run(name)
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

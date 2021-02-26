package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/user"
	"regexp"
	"strings"
	"time"

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

	ctx, err := conf.GetCurrentContext()
	if err != nil {
		fmt.Println(err)
		return repo, err
	}

	if ctx.StoreLocation != nil && strings.HasPrefix(*ctx.StoreLocation, "http") {
		repo = persistance.ApiRepository{*ctx.StoreLocation, *ctx.AuthToken}
	} else if ctx.StoreLocation != nil && strings.HasPrefix(*ctx.StoreLocation, "postgres") {
		re := regexp.MustCompile(`(.*):\/\/(.*):(.*)\@(.*):(.*)\/(.*)`)
		matches := re.FindAllStringSubmatch(*ctx.StoreLocation, -1)

		if len(matches[0]) < 7 {
			return nil, fmt.Errorf("sql dsn expects format postgres://user:password@host:port/database")
		}

		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s database=%s sslmode=disable",
			matches[0][4],
			matches[0][5],
			matches[0][2],
			matches[0][3],
			matches[0][6])

		return persistance.NewSqlRepository(matches[0][1], dsn)
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
	case "-i":
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
	case "-a":
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

	case "-g":
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
		password := generatePassword()

		err = usecases.AddEntry{repo, enc}.Run(name, string(password))
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(password)
	case "-r":
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
	case "-c":
		conf, err := config.LoadCliConfig()
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(os.Args) == 3 {
			conf.CurrentContext = os.Args[2]
			config.SaveCliConfig(conf)
			return
		}

		for _, ctx := range conf.Contexts {
			marker := "  "

			if ctx.Name == conf.CurrentContext {
				marker = "->"
			}

			fmt.Printf("%s %s\n", marker, ctx.Name)
		}
	case "start-server":
		apiConfig, err := config.LoadApiConfig()
		if err != nil {
			fmt.Println(err)
			return
		}

		repo, err := repo()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(api.Api{apiConfig, repo, encryption.NoEncrypter{}, encryption.NoDecrypter{}}.Start())
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

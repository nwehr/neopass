package add

import (
	"fmt"
	"os"

	"github.com/nwehr/neopass"
	"github.com/nwehr/neopass/pkg/config"
	enc "github.com/nwehr/neopass/pkg/encryption/age"
	"golang.org/x/crypto/ssh/terminal"
)

type AddOptions struct {
	What        string
	GetPassword func() (string, error)
	Repo        neopass.EntryRepo
	Encrypter   enc.AgeEncrypter
}

func GetAddOptions(args []string) (AddOptions, error) {
	configOpts, err := config.GetConfigOptions(args)
	if err != nil {
		return AddOptions{}, err
	}

	store, err := configOpts.Config.GetCurrentStore()
	if err != nil {
		return AddOptions{}, err
	}

	repo, err := configOpts.Config.GetCurrentRepo()
	if err != nil {
		return AddOptions{}, err
	}

	enc, err := enc.NewAgeEncrypter(store.Age.Recipients)
	if err != nil {
		return AddOptions{}, err
	}

	return AddOptions{
		What:        args[2],
		GetPassword: ttyPassword,
		Repo:        repo,
		Encrypter:   enc,
	}, nil
}

func RunAdd(opts AddOptions) error {
	plain, err := opts.GetPassword()
	if err != nil {
		return fmt.Errorf("could not get password: %v\n", err)
	}

	encrypted, err := opts.Encrypter.Encrypt(plain)
	if err != nil {
		return fmt.Errorf("could not encrypt password: %v\n", err)
	}

	entry := neopass.Entry{
		Name:     opts.What,
		Password: encrypted,
	}

	if err := opts.Repo.SetEntry(entry); err != nil {
		return fmt.Errorf("could not add entry: %v\n", err)
	}
	return nil
}

func ttyPassword() (string, error) {
	fmt.Print("password: ")

	tty, err := os.Open("/dev/tty")
	if err != nil {
		return "", err
	}

	defer tty.Close()
	defer fmt.Println()

	password, err := terminal.ReadPassword(int(tty.Fd()))
	return string(password), err
}

package add

import (
	"fmt"
	"os"

	"github.com/nwehr/npass"
	"github.com/nwehr/npass/pkg/config"
	enc "github.com/nwehr/npass/pkg/encryption/age"
	"golang.org/x/crypto/ssh/terminal"
)

type AddOptions struct {
	What          string
	ConfigOptions config.ConfigOptions
}

func GetAddOptions(args []string) (AddOptions, error) {
	configOpts, err := config.GetConfigOptions(args)
	if err != nil {
		return AddOptions{}, err
	}

	return AddOptions{
		What:          args[2],
		ConfigOptions: configOpts,
	}, nil
}

func RunAdd(opts AddOptions) error {
	store, err := opts.ConfigOptions.Config.GetCurrentStore()
	if err != nil {
		return fmt.Errorf("could not get current store: %v\n", err)
	}

	r, _ := opts.ConfigOptions.Config.GetCurrentRepo()

	plain, err := ttyPassword()
	if err != nil {
		return fmt.Errorf("could not get password: %v\n", err)
	}

	enc, err := enc.NewAgeEncrypter(store.Age.Recipients)
	if err != nil {
		return fmt.Errorf("could not get encrypter: %v\n", err)
	}

	encrypted, err := enc.Encrypt(string(plain))
	if err != nil {
		return fmt.Errorf("could not encrypt password: %v\n", err)
	}

	entry := npass.Entry{
		Name:     opts.What,
		Password: encrypted,
	}

	if err := r.AddEntry(entry); err != nil {
		return fmt.Errorf("could not add entry: %v\n", err)
	}
	return nil
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

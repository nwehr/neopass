package gen

import (
	"fmt"
	"time"

	mrand "math/rand"

	"github.com/nwehr/neopass/pkg/config"
	enc "github.com/nwehr/neopass/pkg/encryption/age"

	"github.com/atotto/clipboard"
	"github.com/nwehr/neopass"
)

type GenOptions struct {
	Name          string
	ConfigOptions config.ConfigOptions
}

func GetGenOptions(args []string) (GenOptions, error) {
	configOpts, err := config.GetConfigOptions(args)
	if err != nil {
		return GenOptions{}, err
	}

	opts := GenOptions{
		Name:          args[2],
		ConfigOptions: configOpts,
	}

	return opts, nil
}

func RunGen(opts GenOptions) error {
	store, _ := opts.ConfigOptions.Config.GetCurrentStore()
	r, _ := opts.ConfigOptions.Config.GetCurrentRepo()

	plain := genPassword()

	enc, err := enc.NewAgeEncrypter(store.Age.Recipients)
	if err != nil {
		return fmt.Errorf("could not get encrypter: %v\n", err)
	}

	encrypted, err := enc.Encrypt(string(plain))
	if err != nil {
		return fmt.Errorf("could not encrypt password: %v\n", err)
	}

	entry := neopass.Entry{
		Name:     opts.Name,
		Password: encrypted,
	}

	if err := r.SetEntry(entry); err != nil {
		return fmt.Errorf("could not add entry: %v\n", err)
	}

	if err = clipboard.WriteAll(plain); err != nil {
		return fmt.Errorf("coult not write password to clipboard: %v", err)
	}

	fmt.Println("copied to clipboard")
	return nil
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

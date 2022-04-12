package get

import (
	"fmt"

	"github.com/atotto/clipboard"
	"github.com/nwehr/neopass/pkg/config"
	enc "github.com/nwehr/neopass/pkg/encryption/age"
)

type GetOptions struct {
	Name          string
	ConfigOptions config.ConfigOptions
}

func GetGetOptions(args []string) (GetOptions, error) {
	configOpts, err := config.GetConfigOptions(args)
	if err != nil {
		return GetOptions{}, err
	}

	opts := GetOptions{
		Name:          args[1],
		ConfigOptions: configOpts,
	}

	return opts, nil
}

func RunGet(opts GetOptions) error {
	store, err := opts.ConfigOptions.Config.GetCurrentStore()
	if err != nil {
		return err
	}

	r, err := opts.ConfigOptions.Config.GetCurrentRepo()
	if err != nil {
		return err
	}

	entry, err := r.GetEntryByName(opts.Name)
	if err != nil {
		return err
	}

	identity, err := store.Age.UnlockIdentity()
	if err != nil {
		return err
	}

	dec, err := enc.NewAgeDecrypter(identity)
	if err != nil {
		return err
	}

	decrypted, err := dec.Decrypt(entry.Password)
	if err != nil {
		return err
	}

	if err = clipboard.WriteAll(decrypted); err != nil {
		return err
	}

	fmt.Println("copied to clipboard")

	return nil
}

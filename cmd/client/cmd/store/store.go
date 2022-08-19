package store

import (
	"fmt"
	"strings"

	"filippo.io/age"
	"github.com/nwehr/neopass/pkg/config"
)

type StoreOptions struct {
	ShowDetails   bool
	SwitchStore   string
	AddRecipient  string
	ConfigOptions config.ConfigOptions
}

func GetStoreOptions(args []string) (StoreOptions, error) {
	configOpts, err := config.GetConfigOptions(args)
	if err != nil {
		return StoreOptions{}, err
	}

	opts := StoreOptions{
		ConfigOptions: configOpts,
	}

	for i, arg := range args {
		switch arg {
		case "--details":
			opts.ShowDetails = true
		case "--switch":
			opts.SwitchStore = args[i+1]
		case "--add-recipient":
			opts.AddRecipient = args[i+1]
		}
	}

	return opts, nil
}

func RunStore(opts StoreOptions) error {
	if opts.ShowDetails {
		store, err := opts.ConfigOptions.Config.GetCurrentStore()
		if err != nil {
			return err
		}

		_, recipient, err := store.Age.UnlockIdentity()
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println("Name")
		fmt.Println(store.Name)
		fmt.Println()

		fmt.Println("Location")
		fmt.Println(store.Location)
		fmt.Println()

		fmt.Println("Public Identity")
		fmt.Println(recipient)
		fmt.Println()

		fmt.Println("Recipients")
		for _, r := range store.Age.Recipients {
			fmt.Println(r)
		}
		fmt.Println()

		return nil
	}

	if opts.SwitchStore != "" {
		for _, store := range opts.ConfigOptions.Config.Stores {
			if store.Name == opts.SwitchStore {
				opts.ConfigOptions.Config.CurrentStore = opts.SwitchStore
				err := opts.ConfigOptions.Config.WriteFile(config.DefaultConfigFile)
				if err != nil {
					return err
				}

				return nil
			}
		}

		return fmt.Errorf("could not find store '%s'", opts.SwitchStore)
	}

	if opts.AddRecipient != "" {
		index, err := opts.ConfigOptions.Config.GetCurrentStoreIndex()
		if err != nil {
			return err
		}

		_, err = age.ParseRecipients(strings.NewReader(opts.AddRecipient))
		if err != nil {
			return err
		}

		opts.ConfigOptions.Config.Stores[index].Age.Recipients = append(opts.ConfigOptions.Config.Stores[index].Age.Recipients, opts.AddRecipient)

		err = opts.ConfigOptions.Config.WriteFile(config.DefaultConfigFile)
		if err != nil {
			return err
		}

		return nil
	}

	for _, store := range opts.ConfigOptions.Config.Stores {
		marker := "  "

		if store.Name == opts.ConfigOptions.Config.CurrentStore {
			marker = "->"
		}

		fmt.Printf("%s %s\n", marker, store.Name)
	}

	return nil
}

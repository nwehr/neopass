package initstore

import (
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/nwehr/neopass/pkg/config"
)

type InitOptions struct {
	PIV     bool
	PIVSlot string

	NeopassDotCloud bool

	Name string
}

func GetInitOptions(args []string) (InitOptions, error) {
	opts := InitOptions{}

	for i, arg := range args {
		switch arg {
		case "--piv":
			opts.PIV = true
			opts.PIVSlot = getArgValue(args, i)
		case "--neopass.cloud":
			opts.NeopassDotCloud = true

		case "--name":
			opts.Name = args[i+1]
		}
	}

	return opts, nil
}

// TODO: this should be implemented somewhere else
func getArgValue(args []string, argIndex int) string {
	if len(args) > argIndex+1 {
		if args[argIndex+1][0] != '-' {
			return args[argIndex+1]
		}
	}

	return ""
}

func RunInit(opts InitOptions) error {
	newAgeConfig := func() (config.AgeConfig, error) {
		if opts.PIV {
			var slotAddr uint32 = 0x9e

			if opts.PIVSlot != "" {
				addr, err := strconv.ParseUint(opts.PIVSlot, 16, 64)
				if err != nil {
					return config.AgeConfig{}, err
				}

				slotAddr = uint32(addr)
			}

			return config.NewPIVAgeConfig(slotAddr)
		}

		return config.NewDefaultAgeConfig()
	}

	ageConfig, err := newAgeConfig()
	if err != nil {
		return fmt.Errorf("could not setup initial store: %v", err)
	}

	storeConfig := config.StoreConfig{
		Name:     "default",
		Location: config.DefaultStoreFile,
		Age:      ageConfig,
	}

	if opts.NeopassDotCloud {
		storeConfig.Name = "neopass.cloud"
		storeConfig.Location = "https://neopass.cloud?client_uuid=" + uuid.New().String()
	}

	c := config.Config{}
	c.ReadFile(config.DefaultConfigFile)

	if opts.Name != "" {
		storeConfig.Name = opts.Name
	}

	for _, store := range c.Stores {
		if store.Name == storeConfig.Name {
			return fmt.Errorf("%s already exists", store.Name)
		}
	}

	c.Stores = append(c.Stores, storeConfig)
	c.CurrentStore = storeConfig.Name

	err = c.WriteFile(config.DefaultConfigFile)
	if err != nil {
		return fmt.Errorf("could not write default config: %v", err)
	}

	return nil
}

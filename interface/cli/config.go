package cli

import (
	"encoding/json"
	"os"
	"os/user"

	"fmt"

	"github.com/nwehr/npass/infrastructure/encryption/pgp"
	"github.com/nwehr/npass/interface/api"
)

type Config struct {
	CurrentContext string    `json:"currentContext"`
	Contexts       []Context `json:"contexts"`
}

type Context struct {
	Name          string      `json:"name"`
	StoreLocation *string     `json:"store"`
	AuthToken     *string     `json:"authToken"`
	Pgp           *pgp.Config `json:"pgp"`
	Api           *api.Config `json:"api"`
}

func (c Config) GetCurrentContext() (Context, error) {
	for _, ctx := range c.Contexts {
		if c.CurrentContext == ctx.Name {
			return ctx, nil
		}
	}

	return Context{}, fmt.Errorf("No exists")
}

func LoadConfig() (Config, error) {
	usr, _ := user.Current()
	path := usr.HomeDir + "/.npass/npass.json"

	config := Config{}

	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&config)
	return config, err
}

func SaveConfig(config Config) error {
	usr, _ := user.Current()
	path := usr.HomeDir + "/.npass/npass.json"

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer file.Close()

	file.Truncate(0)
	file.Seek(0, 0)

	encoded, err := json.MarshalIndent(config, "", "    ")

	_, err = file.Write(encoded)
	return err
}

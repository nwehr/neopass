package config

import (
	"encoding/json"
	"os"
	"os/user"

	"github.com/nwehr/paws/interface/cli"
)

func LoadCliConfig() (cli.Config, error) {
	usr, _ := user.Current()
	path := usr.HomeDir + "/.paws/paws.json"

	config := cli.Config{}

	file, err := os.Open(path)
	if err != nil {
		return cli.Config{}, err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&config)
	return config, err
}

func SaveCliConfig(config cli.Config) error {
	usr, _ := user.Current()
	path := usr.HomeDir + "/.paws/paws.json"

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

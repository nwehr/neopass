package config

import (
	"encoding/json"
	"os"
	"os/user"

	"github.com/nwehr/paws/infrastructure/encryption/pgp"
)

func LoadPgpConfig() (pgp.Config, error) {
	usr, _ := user.Current()
	path := usr.HomeDir + "/.paws/pgp.json"

	config := pgp.Config{}

	file, err := os.Open(path)
	if err != nil {
		return pgp.Config{}, err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&config)
	return config, err
}

func SavePgpConfig(config pgp.Config) error {
	usr, _ := user.Current()
	path := usr.HomeDir + "/.paws/pgp.json"

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

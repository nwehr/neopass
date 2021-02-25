package config

import (
	"encoding/json"
	"os"
	"os/user"

	"github.com/nwehr/paws/interface/api"
)

func LoadApiConfig() (api.Config, error) {
	listen := os.Getenv("LISTEN")
	authToken := os.Getenv("AUTH_TOKEN")

	if len(listen) > 0 {
		return api.Config{listen, authToken}, nil
	}

	usr, _ := user.Current()
	path := usr.HomeDir + "/.paws/api.json"

	config := api.Config{}

	file, err := os.Open(path)
	if err != nil {
		return api.Config{}, err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&config)
	return config, err
}

func SaveApiConfig(config api.Config) error {
	usr, _ := user.Current()
	path := usr.HomeDir + "/.paws/api.json"

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

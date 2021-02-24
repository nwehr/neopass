package config

import (
	"encoding/json"
	"os"
	"os/user"

	"github.com/nwehr/paws/interface/api"
)

func LoadApiConfig() (api.Config, error) {
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

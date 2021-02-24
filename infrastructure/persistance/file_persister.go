package persistance

import (
	"encoding/json"
	"os"
	"os/user"

	"github.com/nwehr/paws/core/domain"
)

type FilePersister struct {
	Path string
}

func (r FilePersister) Load() (domain.Store, error) {
	store := domain.Store{}

	file, err := os.Open(r.Path)
	if err != nil {
		return store, err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&store)
	return store, err
}

func (r FilePersister) Save(store domain.Store) error {
	file, err := os.OpenFile(r.Path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	defer file.Close()

	file.Truncate(0)
	file.Seek(0, 0)

	encoded, err := json.MarshalIndent(store, "", "    ")

	_, err = file.Write(encoded)
	return err
}

func DefaultFilePersister() FilePersister {
	usr, _ := user.Current()
	path := usr.HomeDir + "/.paws/store.json"

	return FilePersister{Path: path}
}

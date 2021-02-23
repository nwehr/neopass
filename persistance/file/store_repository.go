package file

import (
	"encoding/json"
	"os"
	"os/user"

	"github.com/nwehr/paws/core"
)

type FileStoreRepository struct {
	Path string
}

func (r FileStoreRepository) Load() (core.Store, error) {
	store := core.Store{}

	file, err := os.Open(r.Path)
	if err != nil {
		return store, err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&store)
	return store, err
}

func (r FileStoreRepository) Save(store core.Store) error {
	file, err := os.OpenFile(r.Path, os.O_RDWR, 0)
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

func DefaultStoreRepository() FileStoreRepository {
	usr, _ := user.Current()
	path := usr.HomeDir + "/.paws/store.json"

	return FileStoreRepository{Path: path}
}

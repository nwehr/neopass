package persistance

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"

	"github.com/nwehr/npass/core/domain"
)

type store struct {
	Entries []domain.Entry `json:"entries"`
}

type FileRepository struct {
	Path string
}

func (r FileRepository) AddEntry(entry domain.Entry) error {
	store, err := r.load()
	if err != nil {
		return err
	}

	store.Entries = append(store.Entries, entry)
	return r.save(store)
}

func (r FileRepository) RemoveEntry(name string) error {
	store, err := r.load()
	if err != nil {
		return err
	}

	for i, current := range store.Entries {
		if current.Name == name {
			store.Entries = append(store.Entries[:i], store.Entries[i+1:]...)
			return r.save(store)
		}
	}

	return fmt.Errorf("Not found")
}

func (r FileRepository) GetEntry(name string) (domain.Entry, error) {
	store, err := r.load()
	if err != nil {
		return domain.Entry{}, err
	}

	for _, entry := range store.Entries {
		if entry.Name == name {
			return entry, nil
		}
	}

	return domain.Entry{}, fmt.Errorf("Not found")
}

func (r FileRepository) GetEntryNames() ([]string, error) {
	store, err := r.load()
	if err != nil {
		return nil, err
	}

	names := []string{}

	for _, entry := range store.Entries {
		names = append(names, entry.Name)
	}

	return names, nil
}

func (r FileRepository) load() (store, error) {
	store := store{}

	file, err := os.Open(r.Path)
	if err != nil {
		return store, err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&store)
	return store, err
}

func (r FileRepository) save(store store) error {
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

func DefaultFileRepository() FileRepository {
	usr, _ := user.Current()
	path := usr.HomeDir + "/.npass/store.json"

	return NewFileRepository(path)
}

func NewFileRepository(path string) FileRepository {
	repo := FileRepository{path}

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		repo.save(store{})
	}

	return repo
}

package persistance

import (
	"encoding/json"
	"os"
	"os/user"

	"github.com/nwehr/paws/core/domain"
)

type FileRepository struct {
	Path string
}

func (r FileRepository) AddEntry(entry domain.Entry) error {
	store, err := r.load()
	if err != nil {
		return err
	}

	if err = store.Entries.Add(entry); err != nil {
		return err
	}

	return r.save(store)
}

func (r FileRepository) RemoveEntry(name string) error {
	store, err := r.load()
	if err != nil {
		return err
	}

	if err = store.Entries.Remove(name); err != nil {
		return err
	}

	return r.save(store)
}

func (r FileRepository) GetEntry(name string) (domain.Entry, error) {
	store, err := r.load()
	if err != nil {
		return domain.Entry{}, err
	}

	return store.Entries.Find(name)
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

func (r FileRepository) load() (domain.Store, error) {
	store := domain.Store{}

	file, err := os.Open(r.Path)
	if err != nil {
		return store, err
	}

	defer file.Close()

	err = json.NewDecoder(file).Decode(&store)
	return store, err
}

func (r FileRepository) save(store domain.Store) error {
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
	path := usr.HomeDir + "/.paws/store.json"

	return FileRepository{path}
}

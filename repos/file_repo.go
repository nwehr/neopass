package repos

import (
	"fmt"
	"os"

	"github.com/nwehr/npass"
	"gopkg.in/yaml.v3"
)

type store struct {
	Entries []npass.Entry `yaml:"entries"`
}

type FileRepo struct {
	Path string
}

func (r FileRepo) AddEntry(entry npass.Entry) error {
	store, err := r.load()
	if err != nil {
		return fmt.Errorf("could not load store: %w", err)
	}

	store.Entries = append(store.Entries, entry)
	return r.save(store)
}

func (r FileRepo) RemoveEntryByName(name string) error {
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

	return fmt.Errorf("not found")
}

func (r FileRepo) GetEntryByName(name string) (npass.Entry, error) {
	store, err := r.load()
	if err != nil {
		return npass.Entry{}, err
	}

	for _, entry := range store.Entries {
		if entry.Name == name {
			return entry, nil
		}
	}

	return npass.Entry{}, fmt.Errorf("not found")
}

func (r FileRepo) ListEntryNames() ([]string, error) {
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
func (r FileRepo) load() (store, error) {
	store := store{}

	file, err := os.Open(r.Path)
	if err != nil {
		return store, fmt.Errorf("could not open file: %w", err)
	}

	defer file.Close()

	err = yaml.NewDecoder(file).Decode(&store)
	if err != nil {
		return store, fmt.Errorf("could not decode file: %w", err)
	}

	return store, nil
}

func (r FileRepo) save(store store) error {
	file, err := os.OpenFile(r.Path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}

	defer file.Close()

	file.Truncate(0)
	file.Seek(0, 0)

	enc := yaml.NewEncoder(file)
	defer enc.Close()

	err = enc.Encode(store)
	if err != nil {
		return fmt.Errorf("could not encode store: %w", err)
	}

	return nil
}

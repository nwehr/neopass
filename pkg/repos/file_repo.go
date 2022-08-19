package repos

import (
	"fmt"

	"github.com/nwehr/neopass"
)

type FileRepo struct {
	Path string
}

func NewFileRepo(path string) (neopass.EntryRepo, error) {
	return FileRepo{Path: path}, nil
}

func (r FileRepo) SetEntry(entry neopass.Entry) error {
	store := neopass.Store{}
	if err := store.ReadFile(r.Path); err != nil {
		return fmt.Errorf("could not load store: %w", err)
	}

	store.Entries = append(store.Entries, entry)

	if err := store.WriteFile(r.Path); err != nil {
		return fmt.Errorf("could not save store: %v", err)
	}

	return nil
}

func (r FileRepo) RemoveEntryByName(name string) error {
	store := neopass.Store{}
	if err := store.ReadFile(r.Path); err != nil {
		return fmt.Errorf("could not load store: %w", err)
	}

	for i, current := range store.Entries {
		if current.Name == name {
			store.Entries = append(store.Entries[:i], store.Entries[i+1:]...)

			if err := store.WriteFile(r.Path); err != nil {
				return fmt.Errorf("could not save store: %v", err)
			}

			return nil
		}
	}

	return nil
}

func (r FileRepo) GetEntryByName(name string) (neopass.Entry, error) {
	store := neopass.Store{}
	if err := store.ReadFile(r.Path); err != nil {
		return neopass.Entry{}, fmt.Errorf("could not load store: %w", err)
	}

	for _, entry := range store.Entries {
		if entry.Name == name {
			return entry, nil
		}
	}

	return neopass.Entry{}, fmt.Errorf("not found")
}

func (r FileRepo) ListEntryNames() ([]string, error) {
	store := neopass.Store{}
	if err := store.ReadFile(r.Path); err != nil {
		return nil, fmt.Errorf("could not load store: %w", err)
	}

	names := []string{}

	for _, entry := range store.Entries {
		names = append(names, entry.Name)
	}

	return names, nil
}

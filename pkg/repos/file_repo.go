package repos

import (
	"fmt"

	"github.com/nwehr/npass"
)

type FileRepo struct {
	Path string
}

func (r FileRepo) AddEntry(entry npass.Entry) error {
	store := npass.Store{}
	if err := store.ReadFile(npass.DefaultStoreFile); err != nil {
		return fmt.Errorf("could not load store: %w", err)
	}

	store.Entries = append(store.Entries, entry)

	if err := store.WriteFile(npass.DefaultStoreFile); err != nil {
		return fmt.Errorf("could not save store: %v", err)
	}

	return nil
}

func (r FileRepo) RemoveEntryByName(name string) error {
	store := npass.Store{}
	if err := store.ReadFile(npass.DefaultStoreFile); err != nil {
		return fmt.Errorf("could not load store: %w", err)
	}

	for i, current := range store.Entries {
		if current.Name == name {
			store.Entries = append(store.Entries[:i], store.Entries[i+1:]...)

			if err := store.WriteFile(npass.DefaultStoreFile); err != nil {
				return fmt.Errorf("could not save store: %v", err)
			}

			return nil
		}
	}

	return nil
}

func (r FileRepo) GetEntryByName(name string) (npass.Entry, error) {
	store := npass.Store{}
	if err := store.ReadFile(npass.DefaultStoreFile); err != nil {
		return npass.Entry{}, fmt.Errorf("could not load store: %w", err)
	}

	for _, entry := range store.Entries {
		if entry.Name == name {
			return entry, nil
		}
	}

	return npass.Entry{}, fmt.Errorf("not found")
}

func (r FileRepo) ListEntryNames() ([]string, error) {
	store := npass.Store{}
	if err := store.ReadFile(npass.DefaultStoreFile); err != nil {
		return nil, fmt.Errorf("could not load store: %w", err)
	}

	names := []string{}

	for _, entry := range store.Entries {
		names = append(names, entry.Name)
	}

	return names, nil
}

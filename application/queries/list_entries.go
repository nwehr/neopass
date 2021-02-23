package queries

import "github.com/nwehr/paws/core"

type ListEntries struct {
}

func (q ListEntries) Execute(r core.IStoreRepository) ([]string, error) {
	store, err := r.Load()
	if err != nil {
		return nil, err
	}

	names := []string{}

	for _, entry := range store.Entries {
		names = append(names, entry.Name)
	}

	return names, err
}

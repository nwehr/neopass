package queries

import "github.com/nwehr/paws/core"

type AllEntryNames struct {
}

func (q AllEntryNames) Execute(p core.StorePersister) ([]string, error) {
	store, err := p.Load()
	if err != nil {
		return nil, err
	}

	names := []string{}

	for _, entry := range store.Entries {
		names = append(names, entry.Name)
	}

	return names, err
}

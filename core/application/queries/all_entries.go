package queries

import "github.com/nwehr/paws/core/domain"

type AllEntryNames struct {
}

func (q AllEntryNames) Execute(p domain.StorePersister) ([]string, error) {
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

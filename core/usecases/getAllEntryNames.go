package usecases

import "github.com/nwehr/paws/core/domain"

type GetAllEntryNames struct {
	Persister domain.StorePersister
}

func (u GetAllEntryNames) Run() ([]string, error) {
	store, err := u.Persister.Load()
	if err != nil {
		return nil, err
	}

	names := []string{}

	for _, entry := range store.Entries {
		names = append(names, entry.Name)
	}

	return names, err
}

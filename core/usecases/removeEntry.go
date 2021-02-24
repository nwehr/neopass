package usecases

import "github.com/nwehr/paws/core/domain"

type RemoveEntry struct {
	Persister domain.StorePersister
}

func (u RemoveEntry) Run(name string) error {
	store, err := u.Persister.Load()
	if err != nil {
		return err
	}

	if err = store.Entries.Remove(name); err != nil {
		return err
	}

	return u.Persister.Save(store)
}

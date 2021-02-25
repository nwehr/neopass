package usecases

import "github.com/nwehr/paws/core/domain"

type RemoveEntry struct {
	Repository domain.StoreRepository
}

func (u RemoveEntry) Run(name string) error {
	return u.Repository.RemoveEntry(name)
}

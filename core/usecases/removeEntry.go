package usecases

import "github.com/nwehr/npass/core/domain"

type RemoveEntry struct {
	Repository domain.Repository
}

func (u RemoveEntry) Run(name string) error {
	return u.Repository.RemoveEntry(name)
}

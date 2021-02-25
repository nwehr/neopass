package usecases

import "github.com/nwehr/paws/core/domain"

type GetAllEntryNames struct {
	Repository domain.StoreRepository
}

func (u GetAllEntryNames) Run() ([]string, error) {
	return u.Repository.GetEntryNames()
}

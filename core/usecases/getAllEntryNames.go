package usecases

import (
	"sort"

	"github.com/nwehr/paws/core/domain"
)

type GetAllEntryNames struct {
	Repository domain.StoreRepository
}

func (u GetAllEntryNames) Run() ([]string, error) {
	names, err := u.Repository.GetEntryNames()
	sort.Strings(names)
	return names, err
}

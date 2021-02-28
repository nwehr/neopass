package usecases

import (
	"sort"

	"github.com/nwehr/npass/core/domain"
)

type GetAllEntryNames struct {
	Repository domain.Repository
}

func (u GetAllEntryNames) Run() ([]string, error) {
	names, err := u.Repository.GetEntryNames()
	sort.Strings(names)
	return names, err
}

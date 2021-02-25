package usecases

import (
	"github.com/nwehr/paws/core/domain"
	"github.com/nwehr/paws/infrastructure/encryption"
)

type GetDecryptedPassword struct {
	Repository domain.StoreRepository
	Decrypter  encryption.Decrypter
}

func (u GetDecryptedPassword) Run(name string) (string, error) {
	entry, err := u.Repository.GetEntry(name)
	if err != nil {
		return "", err
	}

	return u.Decrypter.Decrypt(entry.Password)
}

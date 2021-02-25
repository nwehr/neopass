package usecases

import (
	"github.com/nwehr/paws/core/domain"
	"github.com/nwehr/paws/infrastructure/encryption"
)

type AddEntry struct {
	Repository domain.StoreRepository
	Encrypter  encryption.Encrypter
}

func (u AddEntry) Run(name string, password string) error {
	encryptedPassword, err := u.Encrypter.Encrypt(password)
	if err != nil {
		return err
	}

	return u.Repository.AddEntry(domain.Entry{name, encryptedPassword})
}

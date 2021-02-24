package usecases

import (
	"github.com/nwehr/paws/core/domain"
	"github.com/nwehr/paws/infrastructure/encryption"
)

type GetDecryptedPassword struct {
	Persister domain.StorePersister
	Decrypter encryption.Decrypter
}

func (u GetDecryptedPassword) Run(name string) (string, error) {
	store, err := u.Persister.Load()
	if err != nil {
		return "", err
	}

	entry, err := store.Entries.Find(name)
	if err != nil {
		return "", err
	}

	return u.Decrypter.Decrypt(entry.Password)
}

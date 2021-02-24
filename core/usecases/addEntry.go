package usecases

import (
	"github.com/nwehr/paws/core/domain"
	"github.com/nwehr/paws/infrastructure/encryption"
)

type AddEntry struct {
	Encrypter encryption.Encrypter
	Persister domain.StorePersister
}

func (u AddEntry) Run(name string, password string) error {
	store, err := u.Persister.Load()
	if err != nil {
		return err
	}

	encryptedPassword, err := u.Encrypter.Encrypt(password)
	if err != nil {
		return err
	}

	if err := store.Entries.Add(domain.Entry{name, encryptedPassword}); err != nil {
		return err
	}

	return u.Persister.Save(store)
}

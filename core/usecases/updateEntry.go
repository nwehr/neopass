package usecases

import (
	"github.com/nwehr/paws/core/domain"
	"github.com/nwehr/paws/infrastructure/encryption"
)

type UpdateEntry struct {
	Persister domain.StorePersister
	Encrypter encryption.Encrypter
}

func (u UpdateEntry) Run(name, password string) error {
	store, err := u.Persister.Load()
	if err != nil {
		return err
	}

	encryptedPassword, err := u.Encrypter.Encrypt(password)
	if err != nil {
		return err
	}

	if err := store.Entries.Update(domain.Entry{name, encryptedPassword}); err != nil {
		return err
	}

	return u.Persister.Save(store)
}

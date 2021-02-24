package commands

import (
	"github.com/nwehr/paws/core/domain"
	"github.com/nwehr/paws/infrastructure/encryption"
)

type AddEntry struct {
	Name     string
	Password string
}

func (c AddEntry) Execute(enc encryption.Encrypter, p domain.StorePersister) error {
	store, err := p.Load()
	if err != nil {
		return err
	}

	encryptedPassword, err := enc.Encrypt(c.Password)
	if err != nil {
		return err
	}

	if err := store.Entries.Add(domain.Entry{c.Name, encryptedPassword}); err != nil {
		return err
	}

	return p.Save(store)
}

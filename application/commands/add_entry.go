package commands

import (
	"github.com/nwehr/paws/core"
	"github.com/nwehr/paws/encryption"
)

type AddEntry struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (c AddEntry) Execute(e encryption.IEncrypter, r core.IStoreRepository) error {
	store, err := r.Load()
	if err != nil {
		return err
	}

	encryptedPassword, err := e.Encrypt(c.Password)
	if err != nil {
		return err
	}

	if err := store.Entries.Add(core.Entry{c.Name, encryptedPassword}); err != nil {
		return err
	}

	return r.Save(store)
}

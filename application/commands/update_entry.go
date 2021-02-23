package commands

import (
	"github.com/nwehr/paws/core"
	"github.com/nwehr/paws/encryption"
)

type UpdateEntry struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (c UpdateEntry) Execute(e encryption.Encrypter, p core.StorePersister) error {
	store, err := p.Load()
	if err != nil {
		return err
	}

	encryptedPassword, err := e.Encrypt(c.Password)
	if err != nil {
		return err
	}

	if err := store.Entries.Update(core.Entry{c.Name, encryptedPassword}); err != nil {
		return err
	}

	return p.Save(store)
}

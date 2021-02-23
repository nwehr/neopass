package queries

import (
	"github.com/nwehr/paws/core"
	"github.com/nwehr/paws/encryption"
)

type GetEntryPassword struct {
	Name string
}

func (q GetEntryPassword) Execute(d encryption.Decrypter, p core.StorePersister) (string, error) {
	store, err := p.Load()
	if err != nil {
		return "", err
	}

	entry, err := store.Entries.Find(q.Name)
	if err != nil {
		return "", err
	}

	return d.Decrypt(entry.Password)
}

package queries

import (
	"github.com/nwehr/paws/core/domain"
	"github.com/nwehr/paws/infrastructure/encryption"
)

type GetEntryPassword struct {
	Name string
}

func (q GetEntryPassword) Execute(d encryption.Decrypter, p domain.StorePersister) (string, error) {
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

package queries

import (
	"fmt"

	"github.com/nwehr/paws/core"
	"github.com/nwehr/paws/encryption"
)

type GetEntry struct {
	Name string
}

func (q GetEntry) Execute(d encryption.IDecrypter, r core.IStoreRepository) (string, error) {
	store, err := r.Load()
	if err != nil {
		return "", err
	}

	for _, entry := range store.Entries {
		if entry.Name == q.Name {
			return d.Decrypt(entry.Password)
		}
	}

	return "", fmt.Errorf("Not found")
}

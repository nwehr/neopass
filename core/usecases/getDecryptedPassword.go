package usecases

import (
	"github.com/nwehr/npass/core/domain"
	"github.com/nwehr/npass/infrastructure/encryption"
)

type GetDecryptedPassword struct {
	Repository domain.Repository
	Decrypter  encryption.Decrypter
}

func (u GetDecryptedPassword) Run(name string) (string, error) {
	entry, err := u.Repository.GetEntry(name)
	if err != nil {
		return "", err
	}

	return u.Decrypter.Decrypt(entry.Password)
}

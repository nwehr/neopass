package usecases

import (
	"github.com/nwehr/npass/core/domain"
	"github.com/nwehr/npass/infrastructure/encryption"
)

type AddEntry struct {
	Repository domain.Repository
	Encrypter  encryption.Encrypter
}

func (u AddEntry) Run(name string, password string) error {
	encryptedPassword, err := u.Encrypter.Encrypt(password)
	if err != nil {
		return err
	}

	return u.Repository.AddEntry(domain.Entry{Name: name, Password: encryptedPassword})
}

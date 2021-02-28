package usecases

import (
	"github.com/nwehr/npass/core/domain"
	"github.com/nwehr/npass/infrastructure/encryption"
)

type UpdateEntry struct {
	Repository domain.Repository
	Encrypter  encryption.Encrypter
}

func (u UpdateEntry) Run(name, password string) error {
	_, err := u.Repository.GetEntry(name)
	if err != nil {
		return err
	}

	err = u.Repository.RemoveEntry(name)
	if err != nil {
		return err
	}

	encryptedPassword, err := u.Encrypter.Encrypt(password)
	if err != nil {
		return err
	}

	return u.Repository.AddEntry(domain.Entry{Name: name, Password: encryptedPassword})
}

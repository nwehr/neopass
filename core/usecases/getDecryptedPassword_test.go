package usecases

import (
	"testing"

	"github.com/nwehr/paws/infrastructure/encryption"
)

func TestGetDecryptedPassword(t *testing.T) {
	p := DefaultMockPersistor()

	u := AddEntry{p, encryption.NoEncrypter{}}
	u.Run("github.com", "secret1")
	u.Run("gitlab.com", "secret2")
	u.Run("bitbucket.com", "secret3")

	password, err := GetDecryptedPassword{p, encryption.NoDecrypter{}}.Run("gitlab.com")
	if err != nil {
		t.Error(err)
	}

	if password != "secret2" {
		t.Errorf("Expected password to be 'secret2'; got '%s'", password)
	}
}

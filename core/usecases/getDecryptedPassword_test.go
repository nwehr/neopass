package usecases

import (
	"testing"

	"github.com/nwehr/npass/infrastructure/encryption"
)

func TestGetDecryptedPassword(t *testing.T) {
	r := DefaultMockRepository()

	u := AddEntry{r, encryption.NoEncrypter{}}
	u.Run("github.com", "secret1")
	u.Run("gitlab.com", "secret2")
	u.Run("bitbucket.com", "secret3")

	password, err := GetDecryptedPassword{r, encryption.NoDecrypter{}}.Run("gitlab.com")
	if err != nil {
		t.Error(err)
	}

	if password != "secret2" {
		t.Errorf("Expected password to be 'secret2'; got '%s'", password)
	}
}

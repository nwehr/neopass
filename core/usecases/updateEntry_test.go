package usecases

import (
	"testing"

	"github.com/nwehr/npass/infrastructure/encryption"
)

func TestUpdateEntry(t *testing.T) {
	r := DefaultMockRepository()

	u := AddEntry{r, encryption.NoEncrypter{}}
	u.Run("github.com", "abc123")
	u.Run("gitlab.com", "abc123")
	u.Run("bitbucket.com", "abc123")

	UpdateEntry{r, encryption.NoEncrypter{}}.Run("gitlab.com", "secret")

	entry, _ := u.Repository.GetEntry("gitlab.com")

	if entry.Password != "secret" {
		t.Errorf("Expected password to be 'secret'; got '%s'", entry.Password)
	}
}

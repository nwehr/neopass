package usecases

import (
	"testing"

	"github.com/nwehr/paws/infrastructure/encryption"
)

func TestRemoveEntry(t *testing.T) {
	r := DefaultMockRepository()

	u := AddEntry{r, encryption.NoEncrypter{}}
	u.Run("github.com", "abc123")
	u.Run("gitlab.com", "abc123")
	u.Run("bitbucket.com", "abc123")

	names, _ := r.GetEntryNames()

	if len(names) != 3 {
		t.Errorf("Expected 3 entries; got %d", len(names))
	}

	RemoveEntry{r}.Run("gitlab.com")

	names, _ = r.GetEntryNames()

	if len(names) != 2 {
		t.Errorf("Expected 2 entries; got %d", len(names))
	}
}

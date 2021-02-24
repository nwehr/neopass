package usecases

import (
	"testing"
)

func TestUpdateEntry(t *testing.T) {
	p := DefaultMockPersistor()

	u := AddEntry{MockEncrypter{}, p}
	u.Run("github.com", "abc123")
	u.Run("gitlab.com", "abc123")
	u.Run("bitbucket.com", "abc123")

	UpdateEntry{p, MockEncrypter{}}.Run("gitlab.com", "secret")

	store, _ := p.Load()
	entry, _ := store.Entries.Find("gitlab.com")

	if entry.Password != "secret" {
		t.Errorf("Expected password to be 'secret'; got '%s'", entry.Password)
	}
}

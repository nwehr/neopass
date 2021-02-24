package usecases

import "testing"

func TestRemoveEntry(t *testing.T) {
	p := DefaultMockPersistor()

	u := AddEntry{MockEncrypter{}, p}
	u.Run("github.com", "abc123")
	u.Run("gitlab.com", "abc123")
	u.Run("bitbucket.com", "abc123")

	store, _ := p.Load()

	if len(store.Entries) != 3 {
		t.Errorf("Expected 3 entries; got %d", len(store.Entries))
	}

	RemoveEntry{p}.Run("gitlab.com")

	store, _ = p.Load()

	if len(store.Entries) != 2 {
		t.Errorf("Expected 2 entries; got %d", len(store.Entries))
	}
}

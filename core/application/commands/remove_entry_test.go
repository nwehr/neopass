package commands

import "testing"

func TestRemoveEntry(t *testing.T) {
	p := DefaultInMemoryPersistor()

	AddEntry{"github.com", "abc123"}.Execute(NonEncryptor{}, p)
	AddEntry{"gitlab.com", "abc123"}.Execute(NonEncryptor{}, p)
	AddEntry{"bitbucket.com", "abc123"}.Execute(NonEncryptor{}, p)

	store, _ := p.Load()

	if len(store.Entries) != 3 {
		t.Errorf("Expected 3 entries; got %d", len(store.Entries))
	}

	RemoveEntry{"gitlab.com"}.Execute(p)

	store, _ = p.Load()

	if len(store.Entries) != 2 {
		t.Errorf("Expected 2 entries; got %d", len(store.Entries))
	}
}

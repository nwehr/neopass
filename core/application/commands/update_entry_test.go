package commands

import "testing"

func TestUpdateEntry(t *testing.T) {
	p := DefaultInMemoryPersistor()

	AddEntry{"github.com", "abc123"}.Execute(NonEncryptor{}, p)
	AddEntry{"gitlab.com", "abc123"}.Execute(NonEncryptor{}, p)
	AddEntry{"bitbucket.com", "abc123"}.Execute(NonEncryptor{}, p)

	UpdateEntry{"gitlab.com", "secret"}.Execute(NonEncryptor{}, p)

	store, _ := p.Load()
	entry, _ := store.Entries.Find("gitlab.com")

	if entry.Password != "secret" {
		t.Errorf("Expected password to be 'secret'; got '%s'", entry.Password)
	}
}

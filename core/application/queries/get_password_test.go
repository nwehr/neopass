package queries

import (
	"testing"

	"github.com/nwehr/paws/core/application/commands"
)

func TestGetPassword(t *testing.T) {
	p := DefaultInMemoryPersistor()

	commands.AddEntry{"github.com", "secret1"}.Execute(NonEncrypter{}, p)
	commands.AddEntry{"gitlab.com", "secret2"}.Execute(NonEncrypter{}, p)
	commands.AddEntry{"bitbucket.com", "secret3"}.Execute(NonEncrypter{}, p)

	password, err := GetEntryPassword{"gitlab.com"}.Execute(NonDecrypter{}, p)
	if err != nil {
		t.Error(err)
	}

	if password != "secret2" {
		t.Errorf("Expected password to be 'secret2'; got '%s'", password)
	}
}

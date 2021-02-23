package core

import (
	"testing"
)

func TestAddEntry(t *testing.T) {
	s := Store{Identity: "nathan"}
	err := s.Entries.Add(Entry{
		Name:     "github.com",
		Password: "secret",
	})

	if err != nil {
		t.Errorf("Expected err to be nil; got %s", err.Error())
	}

	if len(s.Entries) != 1 {
		t.Errorf("Expected 1 entry; got %d", len(s.Entries))
	}

	err = s.Entries.Add(Entry{
		Name:     "github.com",
		Password: "secret",
	})

	if err == nil {
		t.Errorf("Expected error 'github.com already exists'")
	}

}

func TestUpdateEntry(t *testing.T) {
	s := Store{Identity: "nathan"}

	err := s.Entries.Update(Entry{
		Name:     "github.com",
		Password: "secret2",
	})

	if err == nil {
		t.Errorf("Expected error 'github.com does not exists'")
	}

	s.Entries.Add(Entry{
		Name:     "github.com",
		Password: "secret",
	})

	err = s.Entries.Update(Entry{
		Name:     "github.com",
		Password: "secret2",
	})

	if err != nil {
		t.Errorf("Expected err to be nil")
	}

	if s.Entries[0].Password != "secret2" {
		t.Errorf("Expected password to be 'secret2'; got '%s'", s.Entries[0].Password)
	}

}

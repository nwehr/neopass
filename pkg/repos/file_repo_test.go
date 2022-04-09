package repos

import (
	"testing"

	"github.com/nwehr/neopass"
)

func TestAddEntry(t *testing.T) {
	r := FileRepo{Path: "../store.yaml"}

	{
		if err := r.AddEntry(neopass.Entry{Name: "example.com", Password: "abc123"}); err != nil {
			t.Error(err)
		}
		if err := r.AddEntry(neopass.Entry{Name: "example.net", Password: "123abc"}); err != nil {
			t.Error(err)
		}
	}

	{
		names, err := r.ListEntryNames()
		if err != nil {
			t.Errorf("could not list entry names: %v", err)
		}

		if len(names) != 2 {
			t.Errorf("expected len names to be 2; got %d", len(names))
		}

		if names[0] != "example.com" || names[1] != "example.net" {
			t.Errorf("unexpected entry names: %v", names)
		}

		for _, name := range names {
			if err := r.RemoveEntryByName(name); err != nil {
				t.Errorf("failed to remove entry: %v", err)
			}
		}

		names, err = r.ListEntryNames()
		if err != nil {
			t.Errorf("could not list entry names: %v", err)
		}

		if len(names) != 0 {
			t.Errorf("expected len names to be 0; got %d", len(names))
		}
	}
}

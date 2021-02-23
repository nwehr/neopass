package commands

import (
	"github.com/nwehr/paws/core"
)

type RemoveEntry struct {
	Name string
}

func (c RemoveEntry) Execute(p core.StorePersister) error {
	store, err := p.Load()
	if err != nil {
		return err
	}

	if err = store.Entries.Remove(c.Name); err != nil {
		return err
	}

	return p.Save(store)
}

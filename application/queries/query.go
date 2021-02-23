package queries

import "github.com/nwehr/paws/core"

type IEntryQuery interface {
	Execute(core.IStoreRepository) (string, error)
}

type IEntriesQuery interface {
	Execute(core.IStoreRepository) ([]string, error)
}

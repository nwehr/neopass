package commands

import "github.com/nwehr/paws/core"

type IEntryCommand interface {
	Execute(core.IStoreRepository) error
}

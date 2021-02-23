package core

// Store represents all the password entries for a particular identity
type Store struct {
	Identity string  `json:"identity"`
	Entries  Entries `json:"entries"`
}

type StorePersister interface {
	Load() (Store, error)
	Save(Store) error
}

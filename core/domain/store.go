package domain

// Store represents all the password entries for a particular identity
type Store struct {
	Identity string  `json:"identity"`
	Entries  Entries `json:"entries"`
}

type StorePersister interface {
	Load() (Store, error)
	Save(Store) error
}

type StoreRepository interface {
	AddEntry(Entry) error
	RemoveEntry(string) error
	GetEntry(string) (Entry, error)
	GetEntryNames() ([]string, error)
}

package domain

type Entry struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Entries []Entry

type Repository interface {
	AddEntry(Entry) error
	RemoveEntry(string) error
	GetEntry(string) (Entry, error)
	GetEntryNames() ([]string, error)
}

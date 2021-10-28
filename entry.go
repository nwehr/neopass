package npass

type Entry struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

type EntryRepo interface {
	AddEntry(Entry) error
	RemoveEntryByName(string) error
	GetEntryByName(string) (Entry, error)
	ListEntryNames() ([]string, error)
}

package npass

import "io"

type Entry struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password"`
}

type Entries []Entry

type EntryRepo interface {
	AddEntry(Entry) error
	RemoveEntryByName(string) error
	GetEntryByName(string) (Entry, error)
	ListEntryNames() ([]string, error)
}

type EntryFile struct {
	Entries Entries `yaml:"entries"`
}

func (f *EntryFile) Read(r io.Reader) error {
	return nil
}

func (f EntryFile) Write(w io.Writer) error {
	return nil
}

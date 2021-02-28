package usecases

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/nwehr/npass/core/domain"
	"github.com/nwehr/npass/infrastructure/encryption"
)

func TestAddEntry(t *testing.T) {
	r := DefaultMockRepository()

	u := AddEntry{r, encryption.NoEncrypter{}}
	u.Run("github.com", "abc123")
	u.Run("gitlab.com", "abc123")
	u.Run("bitbucket.com", "abc123")

	names, _ := r.GetEntryNames()

	if len(names) != 3 {
		t.Errorf("Expected 3 entry; got %d", len(names))
	}
}

type store struct {
	Entries domain.Entries `json:"entries"`
}

type MockRepository struct {
	Memory []byte
}

func (p *MockRepository) AddEntry(entry domain.Entry) error {
	store, err := p.load()
	if err != nil {
		return err
	}

	store.Entries = append(store.Entries, entry)

	return p.save(store)
}

func (p *MockRepository) RemoveEntry(name string) error {
	store, err := p.load()
	if err != nil {
		return err
	}

	for i, current := range store.Entries {
		if current.Name == name {
			store.Entries = append(store.Entries[:i], store.Entries[i+1:]...)
			return p.save(store)
		}
	}

	return fmt.Errorf("Not found")
}

func (p MockRepository) GetEntry(name string) (domain.Entry, error) {
	store, err := p.load()
	if err != nil {
		return domain.Entry{}, err
	}

	for _, entry := range store.Entries {
		if entry.Name == name {
			return entry, nil
		}
	}

	return domain.Entry{}, fmt.Errorf("Not found")
}

func (p MockRepository) GetEntryNames() ([]string, error) {
	store, err := p.load()
	if err != nil {
		return nil, err
	}

	names := []string{}

	for _, entry := range store.Entries {
		names = append(names, entry.Name)
	}

	return names, nil
}

func (p MockRepository) load() (store, error) {
	store := store{}
	err := json.Unmarshal(p.Memory, &store)

	return store, err
}

func (p *MockRepository) save(store store) (err error) {
	p.Memory, err = json.Marshal(store)
	return err
}

func DefaultMockRepository() *MockRepository {
	return &MockRepository{[]byte("{\"identity\":\"\", \"entries\": null}")}
}

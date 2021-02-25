package usecases

import (
	"encoding/json"
	"testing"

	"github.com/nwehr/paws/core/domain"
	"github.com/nwehr/paws/infrastructure/encryption"
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

type MockRepository struct {
	Memory []byte
}

func (p *MockRepository) AddEntry(entry domain.Entry) error {
	store, err := p.load()
	if err != nil {
		return err
	}

	if err = store.Entries.Add(entry); err != nil {
		return err
	}

	return p.save(store)
}

func (p *MockRepository) RemoveEntry(name string) error {
	store, err := p.load()
	if err != nil {
		return err
	}

	if err = store.Entries.Remove(name); err != nil {
		return err
	}

	return p.save(store)
}

func (p MockRepository) GetEntry(name string) (domain.Entry, error) {
	store, err := p.load()
	if err != nil {
		return domain.Entry{}, err
	}

	return store.Entries.Find(name)
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

func (p MockRepository) load() (domain.Store, error) {
	store := domain.Store{}
	err := json.Unmarshal(p.Memory, &store)

	return store, err
}

func (p *MockRepository) save(store domain.Store) (err error) {
	p.Memory, err = json.Marshal(store)
	return err
}

func DefaultMockRepository() *MockRepository {
	return &MockRepository{[]byte("{\"identity\":\"\", \"entries\": null}")}
}

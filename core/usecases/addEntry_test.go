package usecases

import (
	"encoding/json"
	"testing"

	"github.com/nwehr/paws/core/domain"
)

func TestAddEntry(t *testing.T) {
	p := DefaultMockPersistor()

	u := AddEntry{MockEncrypter{}, p}
	u.Run("github.com", "abc123")
	u.Run("gitlab.com", "abc123")
	u.Run("bitbucket.com", "abc123")

	store, _ := p.Load()

	if len(store.Entries) != 3 {
		t.Errorf("Expected 3 entry; got %d", len(store.Entries))
	}
}

type MockPersister struct {
	Memory []byte
}

func (p MockPersister) Load() (domain.Store, error) {
	store := domain.Store{}
	err := json.Unmarshal(p.Memory, &store)

	return store, err
}

func (p *MockPersister) Save(store domain.Store) (err error) {
	p.Memory, err = json.Marshal(store)
	return err
}

func DefaultMockPersistor() *MockPersister {
	return &MockPersister{[]byte("{\"identity\":\"\", \"entries\": null}")}
}

type MockEncrypter struct{}

func (MockEncrypter) Encrypt(password string) (string, error) {
	return password, nil
}

type MockDecrypter struct{}

func (MockDecrypter) Decrypt(password string) (string, error) {
	return password, nil
}

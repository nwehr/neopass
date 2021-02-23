package commands

import (
	"encoding/json"
	"testing"

	"github.com/nwehr/paws/core/domain"
)

type InMemoryPersistor struct {
	Memory []byte
}

func (p InMemoryPersistor) Load() (domain.Store, error) {
	store := domain.Store{}
	err := json.Unmarshal(p.Memory, &store)

	return store, err
}

func (p *InMemoryPersistor) Save(store domain.Store) (err error) {
	p.Memory, err = json.Marshal(store)
	return err
}

func DefaultInMemoryPersistor() *InMemoryPersistor {
	return &InMemoryPersistor{[]byte("{\"identity\":\"\", \"entries\": null}")}
}

type NonEncryptor struct{}

func (NonEncryptor) Encrypt(password string) (string, error) {
	return password, nil
}

func TestAddEntry(t *testing.T) {
	p := DefaultInMemoryPersistor()

	AddEntry{"github.com", "abc123"}.Execute(NonEncryptor{}, p)
	AddEntry{"gitlab.com.com", "abc123"}.Execute(NonEncryptor{}, p)
	AddEntry{"bitbucket.com", "abc123"}.Execute(NonEncryptor{}, p)

	store, _ := p.Load()

	if len(store.Entries) != 3 {
		t.Errorf("Expected 3 entry; got %d", len(store.Entries))
	}
}

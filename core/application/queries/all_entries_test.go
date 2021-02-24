package queries

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/nwehr/paws/core/application/commands"
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

type NonEncrypter struct{}

func (NonEncrypter) Encrypt(password string) (string, error) {
	return password, nil
}

type NonDecrypter struct{}

func (NonDecrypter) Decrypt(password string) (string, error) {
	return password, nil
}

func TestAllEntryNames(t *testing.T) {
	p := DefaultInMemoryPersistor()

	commands.AddEntry{"github.com", "abc123"}.Execute(NonEncrypter{}, p)
	commands.AddEntry{"gitlab.com", "abc123"}.Execute(NonEncrypter{}, p)
	commands.AddEntry{"bitbucket.com", "abc123"}.Execute(NonEncrypter{}, p)

	names, err := AllEntryNames{}.Execute(p)
	if err != nil {
		t.Error(err)
	}

	if strings.Join(names, ",") != "github.com,gitlab.com,bitbucket.com" {
		t.Errorf("Expected %s; got %s", "github.com,gitlab.com,bitbucket.com", strings.Join(names, ","))
	}
}

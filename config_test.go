package npass

import (
	"bytes"
	"strings"
	"testing"
)

var fileContents string = `currentStore: default
stores:
  - name: default
    location: /Users/nathanwehr/.npass/local-store.yaml
    age:
	  identity: AGE-SECRET-KEY-162TSCYPCQ23JTZANF2PG489XMYL7CWVAC0U38W2JGY83AAMCXVAQLKYS5R
      recipients:
        - age1995methg04m2h9e6kvetdlgk853ke0u3trzg4gsspz9jc5jhu4nsu82skj
`

func TestReadConfi(t *testing.T) {
	file := strings.NewReader(fileContents)

	c := Config{}
	if err := c.Read(file); err != nil {
		t.Errorf("could not read fileContents: %v", err)
	}

	if c.CurrentStore != "default" {
		t.Errorf("expected default; got %s", c.CurrentStore)
	}

	if c.Stores[0].Name != "default" {
		t.Errorf("expected default; got %s", c.Stores[0].Name)
	}

	if c.Stores[0].Location != "/Users/nathanwehr/.npass/local-store.yaml" {
		t.Errorf("expected /Users/nathanwehr/.npass/local-store.yaml; got %s", c.Stores[0].Location)
	}

	if c.Stores[0].Age.Recipients[0] != "age1995methg04m2h9e6kvetdlgk853ke0u3trzg4gsspz9jc5jhu4nsu82skj" {
		t.Errorf("expected age1995methg04m2h9e6kvetdlgk853ke0u3trzg4gsspz9jc5jhu4nsu82skj; got %s", c.Stores[0].Age.Recipients[0])
	}
}

func TestWriteConfig(t *testing.T) {
	buf := &bytes.Buffer{}

	{
		c := Config{
			CurrentStore: "default",
			Stores: []StoreConfig{
				{
					Name:     "default",
					Location: "~/.npass/store",
				},
			},
		}

		if err := c.Write(buf); err != nil {
			t.Errorf("could not write config: %v", err)
		}
	}

	{
		c := Config{}
		if err := c.Read(buf); err != nil {
			t.Errorf("could not read config: %v", err)
		}

		if c.CurrentStore != "default" {
			t.Errorf("expected default store to be 'default'; got %s", c.CurrentStore)
		}
	}
}

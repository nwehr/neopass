package config

import (
	"bytes"
	"strings"
	"testing"
)

var fileContents string = `currentStore: default
stores:
    - name: default
      location: /Users/nathanwehr/.npass/default-store.yaml
      age:
        identity: AGE-SECRET-KEY-1G6LU2YVE0HW7A45DV43UTY5NTL6P026E3MK9Z454GU9W7USCRXVQ2238CM
        piv:
            slot: 158
        recipients:
            - age163l7rnwpa0ymrn0q0vezmlnhe85l4pvpeksenw3jfmye08ge44dq7audfa`

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

	if c.Stores[0].Location != "/Users/nathanwehr/.npass/default-store.yaml" {
		t.Errorf("expected /Users/nathanwehr/.npass/local-store.yaml; got %s", c.Stores[0].Location)
	}

	if c.Stores[0].Age.Recipients[0] != "age163l7rnwpa0ymrn0q0vezmlnhe85l4pvpeksenw3jfmye08ge44dq7audfa" {
		t.Errorf("expected age163l7rnwpa0ymrn0q0vezmlnhe85l4pvpeksenw3jfmye08ge44dq7audfa; got %s", c.Stores[0].Age.Recipients[0])
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

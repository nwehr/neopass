package usecases

import (
	"strings"
	"testing"
)

func TestGetAllEntryNames(t *testing.T) {
	p := DefaultMockPersistor()

	u := AddEntry{MockEncrypter{}, p}
	u.Run("github.com", "abc123")
	u.Run("gitlab.com", "abc123")
	u.Run("bitbucket.com", "abc123")

	names, err := GetAllEntryNames{p}.Run()
	if err != nil {
		t.Error(err)
	}

	if strings.Join(names, ",") != "github.com,gitlab.com,bitbucket.com" {
		t.Errorf("Expected %s; got %s", "github.com,gitlab.com,bitbucket.com", strings.Join(names, ","))
	}
}

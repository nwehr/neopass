package usecases

import (
	"strings"
	"testing"

	"github.com/nwehr/npass/infrastructure/encryption"
)

func TestGetAllEntryNames(t *testing.T) {
	r := DefaultMockRepository()

	u := AddEntry{r, encryption.NoEncrypter{}}
	u.Run("github.com", "abc123")
	u.Run("gitlab.com", "abc123")
	u.Run("bitbucket.com", "abc123")

	names, _ := r.GetEntryNames()

	if strings.Join(names, ",") != "github.com,gitlab.com,bitbucket.com" {
		t.Errorf("Expected %s; got %s", "github.com,gitlab.com,bitbucket.com", strings.Join(names, ","))
	}
}

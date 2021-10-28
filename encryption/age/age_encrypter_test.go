package age

import (
	"testing"

	"filippo.io/age"
)

func TestEncryptDecrypt(t *testing.T) {
	recip, _ := age.ParseX25519Recipient("age1995methg04m2h9e6kvetdlgk853ke0u3trzg4gsspz9jc5jhu4nsu82skj")
	enc := AgeEncrypter{Recipients: []age.Recipient{recip}}

	encrypted, err := enc.Encrypt("abc123")
	if err != nil {
		t.Error(err)
	}

	ident, _ := age.ParseX25519Identity("AGE-SECRET-KEY-162TSCYPCQ23JTZANF2PG489XMYL7CWVAC0U38W2JGY83AAMCXVAQLKYS5R")
	dec := AgeDecrypter{Identities: []age.Identity{ident}}

	decrypted, err := dec.Decrypt(encrypted)
	if err != nil {
		t.Error(err)
	}

	if decrypted != "abc123" {
		t.Errorf("expected decrypted password to be 'abc123'; got '%s'", decrypted)
	}
}

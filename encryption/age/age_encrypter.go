package age

import (
	"bytes"
	"encoding/ascii85"
	"fmt"
	"io"
	"strings"

	"filippo.io/age"
)

type AgeEncrypter struct {
	Recipients []age.Recipient
}

func NewAgeEncrypter(recipientStrs []string) (AgeEncrypter, error) {
	recipients, err := age.ParseRecipients(strings.NewReader(strings.Join(recipientStrs, "\n")))
	if err != nil {
		return AgeEncrypter{}, err
	}

	return AgeEncrypter{
		Recipients: recipients,
	}, nil
}

func (e AgeEncrypter) Encrypt(password string) (string, error) {
	encoded := new(bytes.Buffer)

	encoder := ascii85.NewEncoder(encoded)
	encrypter, err := age.Encrypt(encoder, e.Recipients...)
	if err != nil {
		return "", fmt.Errorf("could not setup encrypter: %w", err)
	}

	if _, err := io.WriteString(encrypter, password); err != nil {
		return "", fmt.Errorf("could not write password to encrypter: %w", err)
	}

	encrypter.Close()
	encoder.Close()

	return encoded.String(), nil
}

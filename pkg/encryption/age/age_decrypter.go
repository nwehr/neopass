package age

import (
	"bytes"
	"encoding/ascii85"
	"fmt"
	"io"
	"os"
	"strings"

	"filippo.io/age"
	"golang.org/x/crypto/ssh/terminal"
)

type AgeDecrypter struct {
	Identities []age.Identity
}

func NewAgeDecrypter(identity age.Identity) (AgeDecrypter, error) {
	// file, err := os.Open(identityFile)
	// if err != nil {
	// 	return AgeDecrypter{}, fmt.Errorf("could not open identity file: %w", err)
	// }

	// defer file.Close()

	// identities, err := age.ParseIdentities(file)
	// if err != nil {
	// 	err = nil

	// 	password, err := ttyPassword()
	// 	if err != nil {
	// 		return AgeDecrypter{}, fmt.Errorf("could not get password from tty: %w", err)
	// 	}

	// 	passwordIdentity, err := age.NewScryptIdentity(password)
	// 	if err != nil {
	// 		return AgeDecrypter{}, fmt.Errorf("could not get scrypt identity from password: %w", err)
	// 	}

	// 	file.Seek(0, 0)

	// 	rd, err := age.Decrypt(file, passwordIdentity)
	// 	if err != nil {
	// 		return AgeDecrypter{}, fmt.Errorf("could not decrypt identity file: %w", err)
	// 	}

	// 	identities, err = age.ParseIdentities(rd)

	// 	if err != nil {
	// 		return AgeDecrypter{}, fmt.Errorf("could not parse identities: %w", err)
	// 	}
	// }

	return AgeDecrypter{
		Identities: []age.Identity{identity},
	}, nil
}

func (d AgeDecrypter) Decrypt(text string) (string, error) {
	decoder := ascii85.NewDecoder(strings.NewReader(text))
	decrypter, err := age.Decrypt(decoder, d.Identities...)
	if err != nil {
		return "", fmt.Errorf("could not setup decrypter: %w", err)
	}

	decrypted := &bytes.Buffer{}
	if _, err := io.Copy(decrypted, decrypter); err != nil {
		return "", fmt.Errorf("could not copy from decrypter to decrypted: %w", err)
	}

	return decrypted.String(), nil
}

func ttyPassword() (string, error) {
	fmt.Print("password: ")

	tty, err := os.Open("/dev/tty")
	if err != nil {
		return "", err
	}

	defer tty.Close()
	defer fmt.Println()

	password, err := terminal.ReadPassword(int(tty.Fd()))

	return string(password), err
}

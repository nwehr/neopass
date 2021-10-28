package encryption

type Decrypter interface {
	Decrypt(string) (string, error)
}

type NoDecrypter struct{}

func (NoDecrypter) Decrypt(password string) (string, error) {
	return password, nil
}

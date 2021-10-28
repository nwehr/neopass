package encryption

type Encrypter interface {
	Encrypt(string) (string, error)
}

type NoEncrypter struct{}

func (NoEncrypter) Encrypt(password string) (string, error) {
	return password, nil
}

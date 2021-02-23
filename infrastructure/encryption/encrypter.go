package encryption

type Encrypter interface {
	Encrypt(string) (string, error)
}

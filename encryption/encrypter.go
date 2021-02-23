package encryption

type IEncrypter interface {
	Encrypt(string) (string, error)
}

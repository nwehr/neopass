package encryption

type IDecrypter interface {
	Decrypt(string) (string, error)
}

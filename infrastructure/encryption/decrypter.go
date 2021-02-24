package encryption

type Decrypter interface {
	Decrypt(string) (string, error)
}

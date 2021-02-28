package pgp

type Config struct {
	Identities        []string `json:"identities"`
	PublicKeyringPath string   `json:"pubring"`
	SecretKeyringPath string   `json:"secring"`
}

package pgp

type Config struct {
	Identity          string `json:"identity"`
	PublicKeyringPath string `json:"pubring"`
	SecretKeyringPath string `json:"secring"`
}

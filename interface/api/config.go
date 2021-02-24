package api

type Config struct {
	Listen    string `json:"listen"`
	AuthToken string `json:"authToken"`
}

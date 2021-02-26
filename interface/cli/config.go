package cli

import "fmt"

type Config struct {
	CurrentContext string    `json:"currentContext"`
	Contexts       []Context `json:"contexts"`
}

type Context struct {
	Name          string  `json:"name"`
	StoreLocation *string `json:"store"`
	AuthToken     *string `json:"authToken"`
}

func (c Config) GetCurrentContext() (Context, error) {
	for _, ctx := range c.Contexts {
		if c.CurrentContext == ctx.Name {
			return ctx, nil
		}
	}

	return Context{}, fmt.Errorf("No exists")
}

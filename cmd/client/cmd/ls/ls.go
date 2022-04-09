package ls

import (
	"fmt"

	"github.com/nwehr/npass/pkg/config"
)

func RunLs(opts config.ConfigOptions) error {
	r, _ := opts.Config.GetCurrentRepo()
	names, err := r.ListEntryNames()
	if err != nil {
		return fmt.Errorf("could not list entry names: %v\n", err)
	}

	for _, name := range names {
		fmt.Println(name)
	}
	return nil
}

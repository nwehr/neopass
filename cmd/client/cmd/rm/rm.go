package rm

import (
	"fmt"

	"github.com/nwehr/neopass/pkg/config"
)

type RmOptions struct {
	Name          string
	ConfigOptions config.ConfigOptions
}

func GetRmOptions(args []string) (RmOptions, error) {
	configOpts, err := config.GetConfigOptions(args)
	if err != nil {
		return RmOptions{}, err
	}

	opts := RmOptions{
		Name:          args[2],
		ConfigOptions: configOpts,
	}

	return opts, nil
}

func RunRm(opts RmOptions) error {
	r, err := opts.ConfigOptions.Config.GetCurrentRepo()

	if err = r.RemoveEntryByName(opts.Name); err != nil {
		return fmt.Errorf("could not remove entry : %v", err)
	}

	return nil
}

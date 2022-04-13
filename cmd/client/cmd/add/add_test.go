package add

import (
	"testing"

	"github.com/nwehr/expect"
)

func TestGetAddOptions(t *testing.T) {
	args := []string{"neopass", "add", "github.com"}
	opts, err := GetAddOptions(args)

	expect.T(t).NoError(err)
	expect.T(t).String(opts.What).ToEqual(args[2])
}

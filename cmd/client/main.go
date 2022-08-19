package main

import (
	"fmt"
	"os"

	"github.com/nwehr/neopass/cmd/client/cmd/add"
	"github.com/nwehr/neopass/cmd/client/cmd/gen"
	"github.com/nwehr/neopass/cmd/client/cmd/get"
	"github.com/nwehr/neopass/cmd/client/cmd/initstore"
	"github.com/nwehr/neopass/cmd/client/cmd/ls"
	"github.com/nwehr/neopass/cmd/client/cmd/rm"
	"github.com/nwehr/neopass/cmd/client/cmd/store"
	"github.com/nwehr/neopass/pkg/config"
)

var (
	version string
)

func main() {
	if len(os.Args) == 1 {
		showUsage()
		return
	}

	switch os.Args[1] {
	case "init":
		opts, err := initstore.GetInitOptions(os.Args)
		if err != nil {
			Fatalf("could not get init options: %v\n", err)
		}

		err = initstore.RunInit(opts)
		if err != nil {
			Fatalf("%v\n", err)
		}

	case "list":
		opts, err := config.GetConfigOptions(os.Args)
		if err != nil {
			Fatalf("%v\n", err)
		}

		err = ls.RunLs(opts)
		if err != nil {
			Fatalf("%v\n", err)
		}

	case "set":
		opts, err := add.GetAddOptions(os.Args)
		if err != nil {
			Fatalf("%v\n", err)
		}

		err = add.RunAdd(opts)
		if err != nil {
			Fatalf("%v\n", err)
		}

	case "get":
		opts, err := get.GetGetOptions(os.Args[1:])
		if err != nil {
			Fatalf("%v\n", err)
		}

		err = get.RunGet(opts)
		if err != nil {
			Fatalf("%v\n", err)
		}

	case "gen":
		opts, err := gen.GetGenOptions(os.Args)
		if err != nil {
			Fatalf("%v\n", err)
		}

		err = gen.RunGen(opts)
		if err != nil {
			Fatalf("%v\n", err)
		}

	case "rm":
		opts, err := rm.GetRmOptions(os.Args)
		if err != nil {
			Fatalf("%v\n", err)
		}

		err = rm.RunRm(opts)
		if err != nil {
			Fatalf("%v\n", err)
		}

	case "store":
		opts, err := store.GetStoreOptions(os.Args)
		if err != nil {
			Fatalf("%v\n", err)
		}

		err = store.RunStore(opts)
		if err != nil {
			Fatalf("%v\n", err)
		}

	default:
		showUsage()
	}
}

func showUsage() {
	fmt.Printf("neopass %s\n\n", version)

	fmt.Println("Usage")
	fmt.Println("  neopass [<command> [<options>]]")
	fmt.Println()
	fmt.Println("Commands")
	fmt.Println("  list")
	fmt.Println("  set  <name>")
	fmt.Println("  get  <name>")
	fmt.Println("  gen  <name>")
	fmt.Println("  rm   <name>")
	fmt.Println("  init <options>")
	fmt.Println("  store [<options>]")
	fmt.Println()
	fmt.Println("Options")
	fmt.Println("  --config <path>   Use config at path (default ~/.neopass/config.yaml)")
	fmt.Println()
	fmt.Println("Init Options")
	fmt.Println("  --piv [<slot>]        Use PIV card instead of master password")
	fmt.Println("  --name <name>         Specify name for password store")
	fmt.Println("  --neopass.cloud       Use neopass.cloud as password store")
	fmt.Println("  --client-uuid <uuid>  Specify client uuid for neopass.cloud")
	fmt.Println()
	fmt.Println("Store Options")
	fmt.Println("  --switch <name>  Switch to store identified by name")
	fmt.Println("  --details        Show details for current store")
	fmt.Println()
	fmt.Println("Examples")
	fmt.Println("   Initialize new password store on neopass.cloud protected by yubikey")
	fmt.Println("       neopass init --piv --neopass.cloud")
	fmt.Println()
	fmt.Println("   Switch to a password store named default")
	fmt.Println("       neopass store --switch default")
	fmt.Println()
	fmt.Println("Donate")
	fmt.Println("   Bitcoin   (BTC) bc1qkm8gm3ggu8s4lnnc8mp0fahksp23u965hp758c")
	fmt.Println("   Ravencoin (RVN) RSm7jfUjynsVptGyEDaW5yShiXbKBPsHNm")

}

func Fatalf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

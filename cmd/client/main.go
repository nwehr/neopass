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
		opts, err := config.GetConfigOptions(os.Args)
		if err != nil {
			Fatalf("could not load config: %v\n", err)
		}

		if len(os.Args) == 3 {
			for _, store := range opts.Config.Stores {
				if store.Name == os.Args[2] {
					opts.Config.CurrentStore = os.Args[2]
					err = opts.Config.WriteFile(config.DefaultConfigFile)
					if err != nil {
						Fatalf(err.Error())

					}
					return
				}
			}

			Fatalf("could not find store '%s'\n", os.Args[2])
		}

		for _, store := range opts.Config.Stores {
			marker := "  "

			if store.Name == opts.Config.CurrentStore {
				marker = "->"
			}

			fmt.Printf("%s %s\n", marker, store.Name)
		}
	default:
		showUsage()
	}
}

func showUsage() {
	fmt.Printf("neopass %s\n\n", version)

	fmt.Println("Usage")
	fmt.Println("  neopass [<command> [<name>]]")
	fmt.Println()
	fmt.Println("Commands")
	fmt.Println("  list            List entry names")
	fmt.Println("  set   <name>    Set entry identified by name")
	fmt.Println("  get   <name>    Get entry identified by name")
	fmt.Println("  gen   <name>    Generate entry identified by name")
	fmt.Println("  rm    <name>    Remove entry identified by name")
	fmt.Println()
	fmt.Println("  init [--piv [<slot>]] [--neopass.cloud]  Initialize store")
	fmt.Println("  store [--switch <store name>]            Switch to store identified by name or list stores")
	fmt.Println()
	fmt.Println("Options")
	fmt.Println("  --config <path>   Config path (default ~/.neopass/config.yaml)")
	fmt.Println()
	fmt.Println("Examples")
	fmt.Println("   Initialize new password store on neopass cloud")
	fmt.Println("       neopass init --piv --neopass.cloud")
	fmt.Println()
	fmt.Println("   Set an entry for github.com")
	fmt.Println("       neopass set github.com")
	fmt.Println()
	fmt.Println("   Get password for github.com")
	fmt.Println("       neopass get github.com")
	fmt.Println()
	fmt.Println("   Switch to a password store named default")
	fmt.Println("       neopass store --switch default")
}

func Fatalf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

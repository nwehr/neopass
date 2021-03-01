package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/nwehr/npass/core/domain"
	"github.com/nwehr/npass/infrastructure/encryption"
	"github.com/nwehr/npass/infrastructure/encryption/pgp"
	"github.com/nwehr/npass/infrastructure/persistance"
	"github.com/nwehr/npass/interface/api"
	"github.com/nwehr/npass/interface/cli"
)

func main() {
	conf, err := cli.LoadConfig()
	if err != nil {
		fmt.Println("could not open config:", err)
		cli.Cli{nil, nil, nil}.Init()
		return
	}

	ctx, _ := conf.GetCurrentContext()

	repo, err := repo(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(os.Args) > 1 && os.Args[1] == "--server" {
		fmt.Println(api.Api{*ctx.Api, repo, encryption.NoEncrypter{}, encryption.NoDecrypter{}}.Start())
		return
	}

	enc, err := pgp.DefaultEncrypter(*ctx.Pgp)
	if err != nil {
		fmt.Println(err)
	}

	dec, err := pgp.DefaultDecrypter(*ctx.Pgp)
	if err != nil {
		fmt.Println(err)
	}

	cli.Cli{repo, enc, dec}.Start(os.Args)
}

func repo(ctx cli.Context) (domain.Repository, error) {
	if ctx.StoreLocation != nil && strings.HasPrefix(*ctx.StoreLocation, "http") {
		return persistance.ApiRepository{*ctx.StoreLocation, *ctx.AuthToken}, nil
	}

	if ctx.StoreLocation != nil && strings.HasPrefix(*ctx.StoreLocation, "postgres") {
		re := regexp.MustCompile(`(.*):\/\/(.*):(.*)\@(.*):(.*)\/(.*)`)
		matches := re.FindAllStringSubmatch(*ctx.StoreLocation, -1)

		if len(matches[0]) < 7 {
			return nil, fmt.Errorf("sql dsn expects format postgres://user:password@host:port/database")
		}

		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s database=%s sslmode=disable",
			matches[0][4],
			matches[0][5],
			matches[0][2],
			matches[0][3],
			matches[0][6])

		return persistance.NewSqlRepository(matches[0][1], dsn)
	}

	if ctx.StoreLocation != nil {
		return persistance.NewFileRepository(*ctx.StoreLocation), nil
	}

	return persistance.DefaultFileRepository(), nil
}

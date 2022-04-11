tag = $(shell git tag --sort=committerdate | tail -1)
flags = -ldflags "-X main.Version=$(tag)"

all:
	go build $(flags) -o neopass cmd/client/main.go

clean:
	rm neopass

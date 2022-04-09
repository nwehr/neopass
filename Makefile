sha = $(shell git rev-parse --short=8 HEAD)
flags = -ldflags "-X main.Version=$(sha) -X main.Built=$(shell date "+%Y-%m-%d_%H:%M:%S")"

all:
	go build $(flags) -o neopass cmd/client/main.go

clean:
	rm neopass

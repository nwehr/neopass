tag = $(shell git tag --sort=committerdate | tail -1)
flags = -ldflags "-X main.Version=$(tag)"

all:
	go build $(flags) -o neopass cmd/client/main.go

lambda:
	GOOS=linux GOARCH=amd64 go build -o main cmd/lambda/lambda.go
	zip functions.zip main

clean:
	rm neopass

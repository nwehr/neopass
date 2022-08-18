COMMIT = $(shell git rev-parse --short=8 HEAD)
FLAGS = -ldflags "-X 'main.version=${COMMIT}'"
TARGET = neopass

all:
	go build $(FLAGS) -o $(TARGET) cmd/client/main.go

lambda:
	GOOS=linux GOARCH=amd64 go build -o main cmd/lambda/lambda.go
	zip functions.zip main

clean:
	rm -f neopass main functions.zip

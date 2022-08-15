COMMIT = $(shell git rev-parse --short=8 HEAD)
BUILD_DATE := $(shell date '+%Y-%m-%d %H:%M:%S')
FLAGS = -ldflags "-X 'main.commit=${COMMIT}' -X 'main.buildDate=${BUILD_DATE}'"
TARGET = neopass

all:
	go build $(FLAGS) -o $(TARGET) cmd/client/main.go

lambda:
	GOOS=linux GOARCH=amd64 go build -o main cmd/lambda/lambda.go
	zip functions.zip main

clean:
	rm neopass

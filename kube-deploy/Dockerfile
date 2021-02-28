FROM golang:latest

WORKDIR /neopass
COPY . .

RUN go get -d -v ./... \
    && go install -v ./...

CMD ["neopass", "start-server"]
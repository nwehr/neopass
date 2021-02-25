FROM golang:latest

WORKDIR /paws
COPY . .

RUN go get -d -v ./... \
    && go install -v ./...

CMD ["paws", "start-server"]
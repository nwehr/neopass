FROM golang:latest

WORKDIR /npass
COPY . .

RUN go get -d -v ./... \
    && go install -v ./... \
    && mkdir -p /root/.npass

CMD ["npass", "--server"]
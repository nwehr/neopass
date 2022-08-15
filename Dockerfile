FROM golang:1.18
WORKDIR /app
COPY . .
RUN go get ./...
RUN make server

FROM alpine:latest
COPY --from=0 /app/server /app/server
CMD /app/server
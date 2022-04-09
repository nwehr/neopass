#!/bin/sh
go build -o ~/go/bin/neopass ./cmd/client/**.go
cp fzpass.sh ~/go/bin/

#!/bin/sh
go build -o ~/go/bin/npass ./cmd/client/**.go
cp fzpass.sh ~/go/bin/
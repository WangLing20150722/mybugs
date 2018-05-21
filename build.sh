#!/usr/bin/env bash

. ./env.sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mybugs_linux src/main/main_goquery.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o mybugs_win.exe src/main/main_goquery.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o mybugs_darwin src/main/main_goquery.go
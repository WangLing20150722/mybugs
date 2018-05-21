#!/usr/bin/env bash

. ./env.sh
GOOS=linux GOARCH=amd64 go build -o mybugs_linux src/main/main_goquery.go
GOOS=windows GOARCH=amd64 go build -o mybugs_win.exe src/main/main_goquery.go
GOOS=darwin GOARCH=amd64 go build -o mybugs_darwin src/main/main_goquery.go
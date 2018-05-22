#!/usr/bin/env bash

. ./env.sh

set -ex

version=`git log --date=iso --pretty=format:"%cd @%h" -1`
if [ $? -ne 0 ]; then
    version="not a git repo"
fi

compile=`date +"%F %T %z"`" by "`go version`


go build -ldflags "-X \"main.Version=${version}\" -X \"main.Compile=${compile}\"" -o mybugs src/main/main_goquery.go
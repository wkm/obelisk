#!/usr/bin/env sh
go get -u github.com/golang/lint/golint
cp -rv githooks/ .git/hooks/


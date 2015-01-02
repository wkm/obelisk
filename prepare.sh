#!/usr/bin/env sh

echo "Get the linting tool, updating any dependencies"
go get -u github.com/golang/lint/golint

echo "Setup git hooks"
cp -rv githooks/ .git/hooks/

echo "Get and update all dependencies for this project"
go get -t ./...

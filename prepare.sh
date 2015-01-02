#!/usr/bin/env sh

# Get the linting tool, updating any dependencies
go get -u github.com/golang/lint/golint

# Setup git hooks
cp -rv githooks/ .git/hooks/

# Get and update all dependencies for this project
go get -t ./...

#!/bin/bash
clc -s -e eg clip_test.go bin doc.go -L py
go mod tidy
go fmt .
staticcheck .
go vet .
golangci-lint run
unrecognized.py -q
python3 -m flake8 --ignore=W504,W503,E261,E303 .
python3 -m vulture . | grep -v 60%.confidence
git st

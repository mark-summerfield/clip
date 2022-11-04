#!/bin/bash
clc -s -e zOld clop_test.go bin doc.go
go mod tidy
go fmt .
staticcheck .
go vet .
golangci-lint run
git st

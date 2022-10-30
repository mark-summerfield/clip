#!/bin/bash
clc -s -e zOld garg_test.go
go mod tidy
go fmt .
staticcheck .
go vet .
golangci-lint run
git st

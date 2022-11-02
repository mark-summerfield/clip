#!/bin/bash
clc -s -e zOld garg_test.go tester
go mod tidy
go fmt .
staticcheck .
go vet .
golangci-lint run
git st

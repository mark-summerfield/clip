#!/bin/bash
clc -s -e zOld clip_test.go bin doc.go
go mod tidy
go fmt .
staticcheck .
go vet .
golangci-lint run
git st

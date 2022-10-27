#!/bin/bash
clc -s -e zOld
go mod tidy
go fmt .
staticcheck .
go vet .
git st

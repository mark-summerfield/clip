#!/bin/bash
clc -sS -e zOld
go mod tidy
go fmt .
staticcheck .
go vet .
git st

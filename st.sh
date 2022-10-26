#!/bin/bash
clc -sS
go mod tidy
go fmt .
staticcheck .
go vet .
git st

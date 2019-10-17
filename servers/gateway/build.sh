#!/usr/bin/env bash

# go install .
GOOS=linux go build .
# export GOOS="linux go build"
docker build -t jtanderson7/assignment2 .
go clean 

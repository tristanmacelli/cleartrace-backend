#!/usr/bin/env bash

go install .
# export GOOS="linux go build"
docker build -t jtanderson7/assignment2 .
go clean 

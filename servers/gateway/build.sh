#!/bin/bashx

go install .
export GOOS="linux go build"
docker build .
go clean 

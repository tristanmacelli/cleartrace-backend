#!/bin/bashx

go install .
export GOOS="linux"
docker build .
go clean 

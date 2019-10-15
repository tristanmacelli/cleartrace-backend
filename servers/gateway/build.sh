#!/bin/bashx

go install .
docker build .
go clean
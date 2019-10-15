#!/bin/bash

go install .
export GOOS="linux"
docker build .
go clean 
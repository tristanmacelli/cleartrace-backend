#!/usr/bin/env bash

GOOS=linux go build .
docker build -t jtanderson7/assignment2 .
go clean 
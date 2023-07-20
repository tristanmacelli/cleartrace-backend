#!/usr/bin/env bash

echo "Building Go binary & Docker container"
# This creates a go build named after this directory
# CGO_ENABlED=0 will force go to build statically
CGO_ENABLED=0 GOOS=linux GOARCH="amd64" go build .
docker buildx build -t tristanmacelli/gateway . --platform linux/amd64
go clean
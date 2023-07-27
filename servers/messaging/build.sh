#!/usr/bin/env bash

tsc -p ./src --outDir .
docker buildx build -t tristanmacelli/messaging . --platform linux/amd64
rm *.js
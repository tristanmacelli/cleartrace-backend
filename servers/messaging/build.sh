#!/usr/bin/env bash

cd tsc/
tsc --outDir ../
cd -
docker buildx build -t tristanmacelli/messaging . --platform linux/amd64
rm *.js
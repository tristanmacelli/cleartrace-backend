#!/usr/bin/env bash

cd tsc/
tsc --outDir ../
cd -
docker build -t tristanmacelli/messaging .
rm *.js
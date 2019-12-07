#!/usr/bin/env bash

cd tsc/
tsc --outDir ../
cd -
docker build -t jtanderson7/messaging .
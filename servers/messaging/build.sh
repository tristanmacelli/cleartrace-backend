#!/usr/bin/env bash

cd tsc/
tsc --outDir ../
cd -
sudo docker build -t jtanderson7/messaging .
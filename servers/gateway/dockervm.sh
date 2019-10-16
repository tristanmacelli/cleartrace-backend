#!/usr/bin/env bash

docker rm -f jtanderson7/assignment2
docker pull jtanderson7/assignment2
docker run -d -p 80:80 jtanderson7/assignment2
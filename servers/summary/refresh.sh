#!/usr/bin/env bash

export TLSCERT=/etc/letsencrypt/live/api.sauravkharb.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.sauravkharb.me/privkey.pem

docker rm -f summary

# clean up
docker image prune
docker volume prune

docker pull jtanderson7/summary

docker run -d \
-p 5050:5050 \
--network=infrastructure \
--name summary \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
jtanderson7/summary

docker ps
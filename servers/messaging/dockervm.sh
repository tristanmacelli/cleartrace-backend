#!/usr/bin/env bash

docker rm -f messaging
docker rm -f messaging2
docker pull jtanderson7/assignment2

export TLSCERT=/etc/letsencrypt/live/api.sauravkharb.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.sauravkharb.me/privkey.pem

docker rm -f mongodb
docker run --network messagingNetwork -d --name mongodb -v ~/data:/data/db mongo

docker run -d \
--network messagingNetwork \
--name messaging \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e ADDR=5001 \
jtanderson7/assignment2

docker run -d \
--network messagingNetwork \
--name messaging2 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e ADDR=5002 \
jtanderson7/assignment2

docker ps
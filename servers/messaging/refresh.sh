#!/usr/bin/env bash

export TLSCERT=/etc/letsencrypt/live/slack.api.tristanmacelli.com/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/slack.api.tristanmacelli.com/privkey.pem

docker rm -f messaging
docker rm -f messaging2
docker rm -f mongodb

# clean up
docker image prune -f
docker volume prune -f

docker pull tristanmacelli/messaging

docker run -d \
--network=infrastructure \
--name mongodb \
-v ~/data:/data/db \
mongo

docker run -d \
--restart=unless-stopped \
--network=infrastructure \
--name messaging \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e ADDR=5001 \
tristanmacelli/messaging

# docker run -d \
# --restart=unless-stopped \
# --network=infrastructure \
# --name messaging2 \
# -v /etc/letsencrypt:/etc/letsencrypt:ro \
# -e TLSCERT=$TLSCERT \
# -e TLSKEY=$TLSKEY \
# -e ADDR=5002 \
# tristanmacelli/messaging

sleep 1
docker exec -d messaging node startup.js

docker ps
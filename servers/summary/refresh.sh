#!/usr/bin/env bash

export TLSCERT=/etc/letsencrypt/live/slack.api.tristanmacelli.com/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/slack.api.tristanmacelli.com/privkey.pem

docker rm -f summary

# clean up
docker image prune -f
docker volume prune -f

docker pull tristanmacelli/summary

docker run -d \
-p 5050:5050 \
--name summary \
--network=infrastructure \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
tristanmacelli/summary
echo "service refresh completed!"

docker ps
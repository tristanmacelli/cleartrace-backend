#!/usr/bin/env bash

export TLSCERT=/etc/letsencrypt/live/slack.api.tristanmacelli.com/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/slack.api.tristanmacelli.com/privkey.pem

docker rm -f messaging
# docker rm -f messaging2
# docker rm -f userMessageStore

# clean up
docker image prune -f
docker volume prune -f

docker pull tristanmacelli/messaging

# Check that changing the name will change the internal DB name when exec-ing in startup.ts
# docker run -d \
# --network=infrastructure \
# --name userMessageStore \
# -v ~/data:/data/db \
# mongo

docker run -d \
--restart=unless-stopped \
--log-opt max-size=10m \
--log-opt max-file=3 \
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

# sleep 1
# docker exec -t messaging node startup.js

docker ps
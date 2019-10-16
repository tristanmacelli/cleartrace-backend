#!/usr/bin/env bash

docker rm -f jtanderson7/assignment2
docker pull jtanderson7/assignment2

export TLSCERT=/etc/letsencrypt/live/api.sauravkharb.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.sauravkharb.me/privkey.pem

docker run -d \
-p 443:443 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
jtanderson7/assignment2

# docker run -d \
# -p 443:443 \
# -v /etc/letsencrypt/live/api.sauravkharb.me/:/build:ro \
# -e TLSKEY=privkey.pem -e TLSCERT=cert.pem \
# jtanderson7/assignment2;


# etc/letsencrypt/live/api.sauravkharb.me/

# https://api.sauravkharb.me/v1/summary?url=http://ogp.me
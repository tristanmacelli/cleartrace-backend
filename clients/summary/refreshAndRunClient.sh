#!/usr/bin/env bash

docker rm -f jtanderson7/assignment2client
docker pull jtanderson7/assignment2client

export TLSCERT=/etc/letsencrypt/live/a2.sauravkharb.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/a2.sauravkharb.me/privkey.pem

docker run -d \
-p 443:443 \
-p 80:80 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
jtanderson7/assignment2client

docker ps
# docker run -d \
# -p 443:443 \
# -v /etc/letsencrypt/live/api.sauravkharb.me/:/build:ro \
# -e TLSKEY=privkey.pem -e TLSCERT=cert.pem \
# jtanderson7/assignment2;


# etc/letsencrypt/live/api.sauravkharb.me/

# https://api.sauravkharb.me/v1/summary?url=http://ogp.me
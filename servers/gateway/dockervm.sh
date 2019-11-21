#!/usr/bin/env bash

docker rm -f gateway
docker rm -f sessionRedisStore
docker rm -f sqlUserStore
docker pull jtanderson7/assignment2

export TLSCERT=/etc/letsencrypt/live/api.sauravkharb.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.sauravkharb.me/privkey.pem

sudo docker run --name sqlUserStore -d mysql/mysql-server
sudo docker run --name sessionRedisStore -d redis

docker run -d \
-p 443:443 \
--name gateway \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e MESSAGEADDR=messaging:5001,messaging2:5002 \
-e SUMMARYADDR=summary:5050 \
-e SESSIONKEY=sessionkeyrandom \
-e REDISADDR=sessionRedisStore:6379 \
-e DSN=sqlUserStore:3306 \
jtanderson7/assignment2

docker ps
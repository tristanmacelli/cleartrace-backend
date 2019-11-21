#!/usr/bin/env bash

docker rm -f gateway
# TODO: We should probably not be removing the redis & sql on every deploy
docker rm -f userStore
docker rm -f sessionStore
echo "pulling newest version of gateway"
docker pull jtanderson7/assignment2

export TLSCERT=/etc/letsencrypt/live/api.sauravkharb.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.sauravkharb.me/privkey.pem

echo "starting gateway"
docker run --network=infrastructure \
--name userStore -d mysql/mysql-server

docker run --network=infrastructure \
--name sessionStore -d redis

docker run -d \
-p 443:443 \
--network=infrastructure \
--name gateway \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e MESSAGEADDR=messaging:5001,messaging2:5002 \
-e SUMMARYADDR=summary:5050 \
-e SESSIONKEY=sessionkeyrandom \
-e DSN=userStore:3306 \
-e REDISADDR=sessionStore:6379 \
jtanderson7/assignment2
echo "service refresh completed!"

docker ps
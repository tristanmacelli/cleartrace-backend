#!/usr/bin/env bash

export TLSCERT=/etc/letsencrypt/live/slack.api.tristanmacelli.com/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/slack.api.tristanmacelli.com/privkey.pem
export MYSQL_ROOT_PASSWORD=$(openssl rand -base64 18)

docker rm -f gateway
# TODO: We should probably not be removing the redis & sql on every deploy
# docker rm -f userStore # SQL
# docker rm -f sessionStore # REDIS
docker rm -f userMessageQueue

# clean up
echo "cleaning up unused docker artifacts"
docker image prune -f
docker volume prune -f

echo "pulling newest version of gateway"
docker pull tristanmacelli/gateway
docker pull tristanmacelli/db


echo "starting gateway"
# docker run -d --restart=unless-stopped \
# --name userStore \
# --network=infrastructure \
# -e MYSQL_DATABASE=users \
# -e MYSQL_ROOT_PASSWORD=pass \
# -e MYSQL_ROOT_HOST=% \
# tristanmacelli/db

# docker run -d --restart=unless-stopped \
# --name sessionStore \
# --network=infrastructure \
# redis

docker run -d --hostname messagequeue \
--name userMessageQueue \
--network=infrastructure \
rabbitmq:3

sudo docker run -d \
-p 443:443 \
--name gateway \
--network=infrastructure \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e SUMMARYADDR=summary:5050 \
-e MESSAGEADDR=messaging:5001 \
-e SESSIONKEY=sessionkeyrandom \
-e DSN=userStore:3306 \
-e REDISADDR=sessionStore:6379 \
-e MYSQL_ROOT_PASSWORD=pass \
tristanmacelli/gateway
echo "service refresh completed!"

docker ps
#!/usr/bin/env bash

docker rm -f gateway
# TODO: We should probably not be removing the redis & sql on every deploy
# docker rm -f userStore
# docker rm -f sessionStore
docker rm -f rabbitMQ

# clean up
echo "cleaning up unused docker artifacts"
docker image prune -f
docker volume prune -f

echo "pulling newest version of gateway"
docker pull jtanderson7/assignment2
# docker pull jtanderson7/db

export TLSCERT=/etc/letsencrypt/live/api.sauravkharb.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.sauravkharb.me/privkey.pem
export MYSQL_ROOT_PASSWORD=$(openssl rand -base64 18)

echo "starting gateway"
# docker run --restart=unless-stopped \
# --network=infrastructure \
# -e MYSQL_DATABASE=users \
# -e MYSQL_ROOT_PASSWORD=pass \
# -e MYSQL_ROOT_HOST=% \
# --name userStore -d jtanderson7/db 

# docker run --restart=unless-stopped \
# --network=infrastructure \
# --name sessionStore -d redis

docker run -d --network=infrastructure \
--hostname messagequeue --name rabbitMQ rabbitmq:3

sudo docker run -d \
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
-e MYSQL_ROOT_PASSWORD=pass \
jtanderson7/assignment2
echo "service refresh completed!"

docker ps
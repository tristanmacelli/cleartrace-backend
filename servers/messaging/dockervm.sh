#!/usr/bin/env bash

sudo docker rm -f messaging
sudo docker rm -f messaging2
sudo docker pull jtanderson7/assignment2

export TLSCERT=/etc/letsencrypt/live/api.sauravkharb.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.sauravkharb.me/privkey.pem

sudo docker rm -f mongodb
sudo docker run --network=messagingNetwork -d --name mongodb -v ~/data:/data/db mongo
# mongo localhost:27017 collectionCreation.sh (This would create channels & messages)

sudo docker run -d \
--restart=unless-stopped \
--network=messagingNetwork \
--name messaging \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e ADDR=5001 \
jtanderson7/assignment2

sudo docker run -d \
--restart=unless-stopped \
--network=messagingNetwork \
--name messaging2 \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e ADDR=5002 \
jtanderson7/assignment2

docker ps
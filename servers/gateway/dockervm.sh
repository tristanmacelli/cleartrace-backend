#!/usr/bin/env bash

docker rm -f gateway
# TODO: We should probably not be removing the redis & sql on every deploy
docker rm -f userStore
# docker rm -f sessionStore
echo "pulling newest version of gateway"
docker pull jtanderson7/assignment2
docker pull jtanderson7/db

export TLSCERT=/etc/letsencrypt/live/api.sauravkharb.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.sauravkharb.me/privkey.pem
export MYSQL_ROOT_PASSWORD=$(openssl rand -base64 18)

echo "starting gateway"
docker run --restart=unless-stopped \
--network=infrastructure \
-e MYSQL_DATABASE=users \
-e MYSQL_ROOT_PASSWORD=pass \
-e MYSQL_ROOT_HOST=% \
--name userStore -d jtanderson7/db # This was previously mysql (which isn't the container we are creating)

# Create schema for Userstore
# docker run -it \
# --rm \
# -d jtanderson7/db sh -c "mysql -h127.0.0.1 -uroot -p$MYSQL_ROOT_PASSWORD \
#  schema.sql < /docker-entrypoint-initdb.d/schema.sql"
# sudo docker exec userStore sh -c 'exec mysqldump --all-databases -uroot -p"$MYSQL_ROOT_PASSWORD"' > /docker-entrypoint-initdb.d/schema.sql



# docker run --restart=unless-stopped \
# --network=infrastructure \
# --name sessionStore -d redis

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
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
jtanderson7/assignment2
echo "service refresh completed!"

docker ps
  
#!/usr/bin/env bash

docker rm -f messagingClient

# Clean up
docker volume prune
docker image prune

docker pull jtanderson7/messagingClient

export TLSCERT=/etc/letsencrypt/live/a2.sauravkharb.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/a2.sauravkharb.me/privkey.pem

docker run -d \
-p 443:443 \
-p 80:80 \
--name messagingClient \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
jtanderson7/messagingClient

docker ps
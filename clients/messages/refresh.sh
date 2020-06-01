  
#!/usr/bin/env bash

export TLSCERT=/etc/letsencrypt/live/slack.client.tristanmacelli.com/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/slack.client.tristanmacelli.com/privkey.pem

docker rm -f messagingClient

# Clean up
docker volume prune -f
docker image prune -f

docker pull tristanmacelli/messagingclient

docker run -d \
-p 443:443 \
-p 80:80 \
--name messagingClient \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
tristanmacelli/messagingclient

docker ps
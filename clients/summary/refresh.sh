  
#!/usr/bin/env bash

docker rm -f client

# Clean up
docker volume prune -f
docker image prune -f

docker pull jtanderson7/client

export TLSCERT=/etc/letsencrypt/live/a2.sauravkharb.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/a2.sauravkharb.me/privkey.pem

docker run -d \
-p 443:443 \
-p 80:80 \
--name client \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
jtanderson7/client

docker ps
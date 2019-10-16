#!/usr/bin/env bash

./build.sh
docker push jtanderson7/assignment2
# ssh build

chmod g+x ./dockervm.sh
scp -i ~/.ssh/info441_api.pem ./dockervm.sh ec2-user@api.sauravkharb.me:./
ssh -i ~/.ssh/info441_api.pem ec2-user@api.sauravkharb.me "bash ./dockervm.sh"

# cd /etc/letsencrypt/live/api.sauravkharb.me/


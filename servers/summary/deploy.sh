#!/usr/bin/env bash
bash build.sh
docker push jtanderson7/summary
chmod g+x ./refresh.sh
scp -i ~/.ssh/info441_api.pem ./refresh.sh ec2-user@api.sauravkharb.me:./
ssh -i ~/.ssh/info441_api.pem ec2-user@api.sauravkharb.me "bash ./refresh.sh"
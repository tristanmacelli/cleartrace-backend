#!/usr/bin/env bash
cd tsc/
tsc --outDir ../
cd -
bash build.sh
sudo docker push jtanderson7/messaging
chmod g+x ./refresh.sh
sudo scp -i ~/.ssh/info441_api.pem ./refresh.sh ec2-user@api.sauravkharb.me:./
sudo ssh -i ~/.ssh/info441_api.pem ec2-user@api.sauravkharb.me "bash ./refresh.sh"
#!/usr/bin/env bash

echo "build starting..."
bash build.sh
cd ../db
bash buildDb.sh
cd -
echo "build completed!"
docker push jtanderson7/assignment2
docker push jtanderson7/db
chmod g+x ./dockervm.sh
sudo scp -i ~/.ssh/info441_api.pem ./dockervm.sh ec2-user@api.sauravkharb.me:./
echo "service refresh starting..."
sudo ssh -i ~/.ssh/info441_api.pem ec2-user@api.sauravkharb.me "bash ./dockervm.sh"
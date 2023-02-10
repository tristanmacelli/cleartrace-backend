#!/usr/bin/env bash

bash build.sh
cd ../db
# bash buildDb.sh
cd -
echo "build completed!"

echo "Deploying to EC2"
docker push tristanmacelli/gateway
# docker push tristanmacelli/db
chmod g+x ./refresh.sh

echo "Starting Gateway Service.."
ssh -i ~/.ssh/slack-clone ec2-user@slack.api.tristanmacelli.com 'bash -s' < refresh.sh
now=$(date +"%r")
echo "Current time : $now"
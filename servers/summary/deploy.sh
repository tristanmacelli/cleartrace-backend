#!/usr/bin/env bash

bash build.sh

echo "Deploying to EC2"
docker push tristanmacelli/summary
chmod g+x ./refresh.sh

echo "Starting Summary Service"
ssh -i ~/.ssh/slack-clone-server.pem ec2-user@slack.api.tristanmacelli.com 'bash -s' < refresh.sh
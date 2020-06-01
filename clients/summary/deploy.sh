#!/usr/bin/env bash

bash build.sh

echo "Deploying to EC2"
docker push tristanmacelli/summaryclient
chmod g+x ./refresh.sh

echo "Starting Summary Client"
ssh -i ~/.ssh/slack-clone-server.pem ec2-user@slack.client.tristanmacelli.com 'bash -s' < refresh.sh
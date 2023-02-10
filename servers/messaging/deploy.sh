#!/usr/bin/env bash

bash build.sh
docker push tristanmacelli/messaging

ssh -i ~/.ssh/slack-clone ec2-user@slack.api.tristanmacelli.com 'bash -s' < refresh.sh

now=$(date +"%r")
echo "Current time : $now"
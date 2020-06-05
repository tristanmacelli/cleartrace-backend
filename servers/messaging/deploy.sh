#!/usr/bin/env bash

bash build.sh
docker push tristanmacelli/messaging

ssh -i ~/.ssh/slack-clone-server.pem ec2-user@slack.api.tristanmacelli.com 'bash -s' < refresh.sh
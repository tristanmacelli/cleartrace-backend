  
#!/usr/bin/env bash
bash build.sh
docker push tristanmacelli/messagingClient
chmod g+x ./refresh.sh

ssh -i ~/.ssh/slack-clone-server.pem ec2-user@slack.client.tristanmacelli.com < 'bash -s' refresh.sh
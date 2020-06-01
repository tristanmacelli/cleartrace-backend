  
#!/usr/bin/env bash
bash build.sh
docker push tristanmacelli/summaryClient
chmod g+x ./refresh.sh

ssh -i ~/.ssh/slack-clone-client.pem ec2-user@slack.client.tristanmacelli.com < 'bash -s' refresh.sh

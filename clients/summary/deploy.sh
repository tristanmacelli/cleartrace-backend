  
#!/usr/bin/env bash
bash build.sh
docker push jtanderson7/assignment2client
chmod g+x ./refreshAndRunClient.sh
scp -i ~/.ssh/info441_a2.pem ./refreshAndRunClient.sh ec2-user@a2.sauravkharb.me:./
ssh -i ~/.ssh/info441_a2.pem ec2-user@a2.sauravkharb.me "bash ./refreshAndRunClient.sh"

  
#!/usr/bin/env bash
bash build.sh
docker push jtanderson7/assignment2client
chmod g+x ./refresh.sh
scp -i ~/.ssh/info441_a2.pem ./refresh.sh ec2-user@a2.sauravkharb.me:./
ssh -i ~/.ssh/info441_a2.pem ec2-user@a2.sauravkharb.me "bash ./refresh.sh"

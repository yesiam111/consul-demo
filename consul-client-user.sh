#!/bin/bash

svc-name=$1

curl --fail --silent --show-error --location https://apt.releases.hashicorp.com/gpg | \
      gpg --dearmor | \
      sudo dd of=/usr/share/keyrings/hashicorp-archive-keyring.gpg

echo "deb [arch=amd64 signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | \
 sudo tee -a /etc/apt/sources.list.d/hashicorp.list

sudo apt-get update

sudo apt-cache policy consul

sudo apt-get install consul -y

cat <<EOT >> /etc/consul.d/consul.hcl
server = false
encrypt = "qDOPBEr+/oUVeOFQOnVypxwDaHzLrD+lvjo5vCEBbZ0="
retry_join = ["provider=aws tag_key=Name tag_value=consul-server"]
EOT

systemctl enable consul
systemctl start consul

wget https://go.dev/dl/go1.18.1.linux-amd64.tar.gz
tar xzvf go1.18.1.linux-amd64.tar.gz -C /usr/local/
export PATH=$PATH:/usr/local/go/bin

cd consul-demo/user-service
go mod tidy
go run main.go &
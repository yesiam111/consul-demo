#!/bin/bash

curl --fail --silent --show-error --location https://apt.releases.hashicorp.com/gpg | \
      gpg --dearmor | \
      sudo dd of=/usr/share/keyrings/hashicorp-archive-keyring.gpg

echo "deb [arch=amd64 signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(lsb_release -cs) main" | \
 sudo tee -a /etc/apt/sources.list.d/hashicorp.list

sudo apt-get update

sudo apt-cache policy consul

sudo apt-get install consul -y

cat <<EOT >> /etc/consul.d/consul.hcl
server = true
bind_addr = "0.0.0.0"
bootstrap_expect=1
encrypt = "qDOPBEr+/oUVeOFQOnVypxwDaHzLrD+lvjo5vCEBbZ0="
client_addr = "0.0.0.0"
EOT

systemctl enable consul
systemctl start consul
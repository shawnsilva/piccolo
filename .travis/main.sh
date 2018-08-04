#!/bin/bash

set -o errexit

echo "INFO: Installing dependencies."

sudo apt update -y
sudo apt install realpath python python-pip -y
sudo apt install --only-upgrade docker-ce -y

sudo pip install docker-compose || true

docker info
docker-compose --version


echo "INFO: Setup Docker"

echo '{
  "experimental": true,
  "storage-driver": "overlay2",
  "max-concurrent-downloads": 50,
  "max-concurrent-uploads": 50
}' | sudo tee /etc/docker/daemon.json
sudo service docker restart

echo "SUCCESS: Docker Ready"

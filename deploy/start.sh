#!/bin/bash
docker ps -a | grep runbox$ && docker rm -f runbox
docker run -d --name runbox -p 80:8080 --privileged --restart always \
-v /data:/data \
-v /var/run/docker.sock:/var/run/docker.sock \
-v /var/lib/docker/containers:/var/lib/docker/containers \
ghcr.io/zetaoss/runbox

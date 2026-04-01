#!/bin/bash
docker run -d --network host --name host-nginx nginx
docker inspect host-nginx --format '{{.HostConfig.NetworkMode}}' > /tmp/host-network-info.txt

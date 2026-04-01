#!/bin/bash
docker network inspect bridge > /tmp/bridge-info.txt
docker run -d --name bridge-test-1 alpine sleep 3600
docker run -d --name bridge-test-2 alpine sleep 3600
echo "bridge-test-1:" > /tmp/container-ips.txt
docker inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' bridge-test-1 >> /tmp/container-ips.txt
echo "bridge-test-2:" >> /tmp/container-ips.txt
docker inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' bridge-test-2 >> /tmp/container-ips.txt

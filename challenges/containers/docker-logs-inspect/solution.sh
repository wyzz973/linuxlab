#!/bin/bash
docker run -d --name log-demo -p 8082:80 nginx
sleep 2
curl -s http://localhost:8082 > /dev/null 2>&1
docker logs log-demo > /tmp/container-logs.txt 2>&1
docker inspect --format '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' log-demo > /tmp/container-ip.txt

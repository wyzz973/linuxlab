#!/bin/bash
docker run -d --name cp-demo nginx
echo "file from host" > /tmp/upload.txt
docker cp /tmp/upload.txt cp-demo:/tmp/upload.txt
docker cp cp-demo:/etc/nginx/nginx.conf /tmp/nginx.conf

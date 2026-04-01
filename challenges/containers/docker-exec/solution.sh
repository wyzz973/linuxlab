#!/bin/bash
docker run -d --name exec-demo nginx
docker exec exec-demo bash -c 'echo "executed inside container" > /tmp/exec-test.txt'
docker exec exec-demo ps aux > /tmp/container-processes.txt 2>/dev/null || \
docker exec exec-demo ls /proc > /tmp/container-processes.txt

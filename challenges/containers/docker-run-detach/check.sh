#!/bin/bash
if docker ps --format '{{.Names}}' | grep -q '^my-nginx$'; then
    echo "PASS"
    exit 0
fi
echo "FAIL: 名为 my-nginx 的容器未在运行"
exit 1

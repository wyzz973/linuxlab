#!/bin/bash
if docker ps --format '{{.Names}}' | grep -q 'compose-web'; then
    echo "PASS"
    exit 0
fi
echo "FAIL: compose-web 容器未在运行"
exit 1

#!/bin/bash
if ! docker network ls --format '{{.Name}}' | grep -q '^app-network$'; then
    echo "FAIL: 网络 app-network 不存在"
    exit 1
fi
if ! docker ps --format '{{.Names}}' | grep -q '^web-app$'; then
    echo "FAIL: web-app 容器未在运行"
    exit 1
fi
if ! docker ps --format '{{.Names}}' | grep -q '^test-client$'; then
    echo "FAIL: test-client 容器未在运行"
    exit 1
fi
# Check both are on app-network
NET1=$(docker inspect web-app --format '{{range $k,$v := .NetworkSettings.Networks}}{{$k}}{{end}}')
NET2=$(docker inspect test-client --format '{{range $k,$v := .NetworkSettings.Networks}}{{$k}}{{end}}')
if ! echo "$NET1" | grep -q "app-network"; then
    echo "FAIL: web-app 不在 app-network 网络上"
    exit 1
fi
if ! echo "$NET2" | grep -q "app-network"; then
    echo "FAIL: test-client 不在 app-network 网络上"
    exit 1
fi
if [ ! -f /tmp/ping-result.txt ]; then
    echo "FAIL: /tmp/ping-result.txt 不存在"
    exit 1
fi
echo "PASS"
exit 0

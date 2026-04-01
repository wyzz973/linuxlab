#!/bin/bash
for net in frontend backend; do
    if ! docker network ls --format '{{.Name}}' | grep -q "^${net}$"; then
        echo "FAIL: 网络 ${net} 不存在"
        exit 1
    fi
done
if ! docker ps --format '{{.Names}}' | grep -q '^proxy-server$'; then
    echo "FAIL: proxy-server 容器未在运行"
    exit 1
fi
if ! docker ps --format '{{.Names}}' | grep -q '^app-server$'; then
    echo "FAIL: app-server 容器未在运行"
    exit 1
fi
# Check proxy-server is on both networks
NETS=$(docker inspect proxy-server --format '{{range $k,$v := .NetworkSettings.Networks}}{{$k}} {{end}}')
if ! echo "$NETS" | grep -q "frontend"; then
    echo "FAIL: proxy-server 未连接到 frontend 网络"
    exit 1
fi
if ! echo "$NETS" | grep -q "backend"; then
    echo "FAIL: proxy-server 未连接到 backend 网络"
    exit 1
fi
echo "PASS"
exit 0

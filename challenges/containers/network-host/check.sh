#!/bin/bash
if ! docker ps --format '{{.Names}}' | grep -q '^host-nginx$'; then
    echo "FAIL: host-nginx 容器未在运行"
    exit 1
fi
NET_MODE=$(docker inspect host-nginx --format '{{.HostConfig.NetworkMode}}')
if [ "$NET_MODE" != "host" ]; then
    echo "FAIL: 网络模式不是 host，当前为 $NET_MODE"
    exit 1
fi
if [ ! -f /tmp/host-network-info.txt ]; then
    echo "FAIL: /tmp/host-network-info.txt 不存在"
    exit 1
fi
echo "PASS"
exit 0

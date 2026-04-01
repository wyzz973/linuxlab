#!/bin/bash
if ! docker network ls --format '{{.Name}}' | grep -q '^my-network$'; then
    echo "FAIL: 网络 my-network 不存在"
    exit 1
fi
if ! docker network ls --format '{{.Name}}' | grep -q '^isolated-net$'; then
    echo "FAIL: 网络 isolated-net 不存在"
    exit 1
fi
SUBNET=$(docker network inspect isolated-net --format '{{range .IPAM.Config}}{{.Subnet}}{{end}}')
if [ "$SUBNET" != "172.28.0.0/16" ]; then
    echo "FAIL: isolated-net 子网不正确，期望 172.28.0.0/16"
    exit 1
fi
if [ ! -f /tmp/network-list.txt ] || [ ! -f /tmp/network-list-after.txt ]; then
    echo "FAIL: 网络列表文件不存在"
    exit 1
fi
echo "PASS"
exit 0

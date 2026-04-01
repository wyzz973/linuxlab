#!/bin/bash
if ! docker images --format '{{.Repository}}:{{.Tag}}' | grep -q "envapp:v1"; then
    echo "FAIL: 镜像 envapp:v1 未构建"
    exit 1
fi
if ! docker ps --format '{{.Names}}' | grep -q '^env-test$'; then
    echo "FAIL: 容器 env-test 未在运行"
    exit 1
fi
ENV_CHECK=$(docker inspect env-test --format '{{range .Config.Env}}{{println .}}{{end}}')
if ! echo "$ENV_CHECK" | grep -q "APP_NAME=MyApp"; then
    echo "FAIL: 环境变量 APP_NAME 未设置"
    exit 1
fi
if ! echo "$ENV_CHECK" | grep -q "APP_VERSION=1.0"; then
    echo "FAIL: 环境变量 APP_VERSION 未设置"
    exit 1
fi
echo "PASS"
exit 0

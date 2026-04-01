#!/bin/bash
if ! docker ps --format '{{.Names}}' | grep -q 'health-db'; then
    echo "FAIL: health-db 容器未在运行"
    exit 1
fi
if ! docker ps --format '{{.Names}}' | grep -q 'health-web'; then
    echo "FAIL: health-web 容器未在运行"
    exit 1
fi
# Check healthcheck is configured
HEALTH=$(docker inspect health-db --format '{{.Config.Healthcheck}}')
if [ -z "$HEALTH" ] || [ "$HEALTH" = "<nil>" ]; then
    echo "FAIL: health-db 未配置健康检查"
    exit 1
fi
echo "PASS"
exit 0

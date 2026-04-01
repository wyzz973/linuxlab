#!/bin/bash
if [ ! -f /tmp/multi-app/docker-compose.yaml ] && [ ! -f /tmp/multi-app/docker-compose.yml ]; then
    echo "FAIL: docker-compose.yaml 不存在"
    exit 1
fi
COMPOSE_FILE=$(ls /tmp/multi-app/docker-compose.y*ml 2>/dev/null | head -1)
# Check all 3 services are running
RUNNING=$(docker compose -f "$COMPOSE_FILE" ps --format '{{.State}}' 2>/dev/null | grep -c "running")
if [ "$RUNNING" -lt 3 ]; then
    # Fallback: check containers directly
    WEB=$(docker ps --format '{{.Names}}' | grep -c "web")
    DB=$(docker ps --format '{{.Names}}' | grep -c "db\|redis")
    API=$(docker ps --format '{{.Names}}' | grep -c "api\|node")
    TOTAL=$((WEB + DB + API))
    if [ "$TOTAL" -lt 3 ]; then
        echo "FAIL: 未找到3个运行中的服务（找到 $TOTAL 个）"
        exit 1
    fi
fi
echo "PASS"
exit 0

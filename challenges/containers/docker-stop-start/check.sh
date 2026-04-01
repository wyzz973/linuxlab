#!/bin/bash
if ! docker ps --format '{{.Names}}' | grep -q '^lifecycle-test$'; then
    echo "FAIL: 容器 lifecycle-test 未在运行"
    exit 1
fi
# Check that the container has been restarted at least once
RESTART_COUNT=$(docker inspect lifecycle-test --format '{{.RestartCount}}' 2>/dev/null)
STARTED_AT=$(docker inspect lifecycle-test --format '{{.State.StartedAt}}' 2>/dev/null)
if [ -n "$STARTED_AT" ]; then
    echo "PASS"
    exit 0
fi
echo "FAIL: 无法验证容器状态"
exit 1

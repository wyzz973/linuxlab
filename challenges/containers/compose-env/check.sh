#!/bin/bash
if [ ! -f /tmp/compose-env/.env ]; then
    echo "FAIL: .env 文件不存在"
    exit 1
fi
if ! docker ps --format '{{.Names}}' | grep -q 'compose-mysql'; then
    echo "FAIL: compose-mysql 容器未在运行"
    exit 1
fi
ENV_CHECK=$(docker inspect compose-mysql --format '{{range .Config.Env}}{{println .}}{{end}}')
if echo "$ENV_CHECK" | grep -q "MYSQL_ROOT_PASSWORD=rootpass123"; then
    if echo "$ENV_CHECK" | grep -q "MYSQL_DATABASE=mydb"; then
        echo "PASS"
        exit 0
    fi
fi
echo "FAIL: 环境变量未正确设置"
exit 1

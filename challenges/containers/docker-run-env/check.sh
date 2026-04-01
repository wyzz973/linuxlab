#!/bin/bash
if ! docker ps --format '{{.Names}}' | grep -q '^my-mysql$'; then
    echo "FAIL: 容器 my-mysql 未在运行"
    exit 1
fi
ENV_PASS=$(docker inspect my-mysql --format '{{range .Config.Env}}{{println .}}{{end}}' | grep "MYSQL_ROOT_PASSWORD=my-secret-pw")
ENV_DB=$(docker inspect my-mysql --format '{{range .Config.Env}}{{range $i, $e := split .}}{{end}}{{println .}}{{end}}' 2>/dev/null)
if docker inspect my-mysql --format '{{range .Config.Env}}{{println .}}{{end}}' | grep -q "MYSQL_ROOT_PASSWORD=my-secret-pw"; then
    if docker inspect my-mysql --format '{{range .Config.Env}}{{println .}}{{end}}' | grep -q "MYSQL_DATABASE=testdb"; then
        echo "PASS"
        exit 0
    fi
    echo "FAIL: MYSQL_DATABASE 环境变量未设置为 testdb"
    exit 1
fi
echo "FAIL: MYSQL_ROOT_PASSWORD 环境变量未正确设置"
exit 1

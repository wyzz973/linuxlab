#!/bin/bash
if ! docker volume ls --format '{{.Name}}' | grep -q '^app-data$'; then
    echo "FAIL: 数据卷 app-data 不存在"
    exit 1
fi
if ! docker volume ls --format '{{.Name}}' | grep -q '^db-data$'; then
    echo "FAIL: 数据卷 db-data 不存在"
    exit 1
fi
if [ ! -f /tmp/volume-info.txt ]; then
    echo "FAIL: /tmp/volume-info.txt 不存在"
    exit 1
fi
if [ ! -f /tmp/volume-list.txt ]; then
    echo "FAIL: /tmp/volume-list.txt 不存在"
    exit 1
fi
echo "PASS"
exit 0

#!/bin/bash
for tag in "tagdemo:v1" "tagdemo:latest" "myregistry/tagdemo:v1"; do
    if ! docker images --format '{{.Repository}}:{{.Tag}}' | grep -q "^${tag}$"; then
        echo "FAIL: 镜像 ${tag} 不存在"
        exit 1
    fi
done
if [ ! -f /tmp/tag-list.txt ]; then
    echo "FAIL: /tmp/tag-list.txt 不存在"
    exit 1
fi
echo "PASS"
exit 0

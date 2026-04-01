#!/bin/bash
if [ ! -f /tmp/volume-backup.tar ]; then
    echo "FAIL: /tmp/volume-backup.tar 不存在"
    exit 1
fi
if ! docker volume ls --format '{{.Name}}' | grep -q '^backup-restore$'; then
    echo "FAIL: 数据卷 backup-restore 不存在"
    exit 1
fi
# Check that restored data matches
ORIGINAL=$(docker run --rm -v backup-source:/data alpine ls /data 2>/dev/null)
RESTORED=$(docker run --rm -v backup-restore:/data alpine ls /data 2>/dev/null)
if [ -z "$RESTORED" ]; then
    echo "FAIL: 恢复的卷中没有数据"
    exit 1
fi
if [ "$ORIGINAL" = "$RESTORED" ]; then
    echo "PASS"
    exit 0
fi
echo "FAIL: 恢复的数据与原始数据不一致"
exit 1

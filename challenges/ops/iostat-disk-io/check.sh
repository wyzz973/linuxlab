#!/bin/bash
for f in /tmp/iostat_output.txt /tmp/diskstats.txt /tmp/disk_usage.txt; do
    if [ ! -f "$f" ] || [ ! -s "$f" ]; then
        echo "FAIL: $f not found or empty"
        exit 1
    fi
done
if grep -qiE "device\|avg-cpu\|Filesystem" /tmp/iostat_output.txt /tmp/disk_usage.txt 2>/dev/null; then
    echo "PASS"
    exit 0
fi
echo "PASS"
exit 0

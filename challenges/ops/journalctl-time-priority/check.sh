#!/bin/bash
for f in /tmp/today_logs.txt /tmp/error_logs.txt /tmp/hour_logs.txt; do
    if [ ! -f "$f" ]; then
        echo "FAIL: $f not found"
        exit 1
    fi
done
echo "PASS"
exit 0

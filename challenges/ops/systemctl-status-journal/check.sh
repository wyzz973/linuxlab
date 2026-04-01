#!/bin/bash
for f in /tmp/service_detail.txt /tmp/service_logs.txt /tmp/service_errors.txt; do
    if [ ! -f "$f" ]; then
        echo "FAIL: $f not found"
        exit 1
    fi
done
echo "PASS"
exit 0

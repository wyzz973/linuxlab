#!/bin/bash
for f in /tmp/lsof_listen.txt /tmp/lsof_pid.txt; do
    if [ ! -f "$f" ]; then
        echo "FAIL: $f not found"
        exit 1
    fi
done
if [ ! -f /tmp/lsof_deleted.txt ]; then
    echo "FAIL: /tmp/lsof_deleted.txt not found"
    exit 1
fi
if [ -s /tmp/lsof_listen.txt ]; then
    echo "PASS"
    exit 0
fi
echo "FAIL: lsof output is empty"
exit 1

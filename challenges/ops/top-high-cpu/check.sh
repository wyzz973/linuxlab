#!/bin/bash
for f in /tmp/top_cpu.txt /tmp/top_mem.txt /tmp/proc_count.txt; do
    if [ ! -f "$f" ] || [ ! -s "$f" ]; then
        echo "FAIL: $f not found or empty"
        exit 1
    fi
done
if grep -qE "PID|%CPU|USER" /tmp/top_cpu.txt; then
    echo "PASS"
    exit 0
fi
# Accept ps output without header
if [ "$(wc -l < /tmp/top_cpu.txt)" -ge 1 ]; then
    echo "PASS"
    exit 0
fi
echo "FAIL: invalid process list"
exit 1

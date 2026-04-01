#!/bin/bash
for f in /tmp/sar_cpu.txt /tmp/sar_mem.txt /tmp/sar_net.txt; do
    if [ ! -f "$f" ]; then
        echo "FAIL: $f not found"
        exit 1
    fi
done
if [ -s /tmp/sar_cpu.txt ]; then
    echo "PASS"
    exit 0
fi
echo "FAIL: sar output is empty"
exit 1

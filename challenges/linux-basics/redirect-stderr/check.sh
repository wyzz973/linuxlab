#!/bin/bash
if [ ! -f /tmp/stdout.txt ] || [ ! -f /tmp/stderr.txt ]; then
    echo "FAIL: Output files not found"
    exit 1
fi
if [ -s /tmp/stderr.txt ]; then
    if grep -qi "no such file\|cannot access\|not found" /tmp/stderr.txt; then
        echo "PASS"
        exit 0
    fi
fi
echo "FAIL: stderr should contain error about nonexistent directory"
exit 1

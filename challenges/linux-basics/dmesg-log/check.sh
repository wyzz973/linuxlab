#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
lines=$(wc -l < /tmp/result.txt | tr -d ' ')
if [ "$lines" -le 20 ] && [ "$lines" -gt 0 ]; then
    echo "PASS"
    exit 0
elif [ "$lines" -eq 0 ] && [ -f /tmp/result.txt ]; then
    echo "PASS: File created but dmesg may be empty in container"
    exit 0
else
    echo "FAIL: Expected at most 20 lines, got $lines"
    exit 1
fi

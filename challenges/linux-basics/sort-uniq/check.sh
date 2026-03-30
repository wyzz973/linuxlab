#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
# First line should have the highest count (192.168.1.100 with 5)
first=$(head -1 /tmp/result.txt | tr -s ' ')
if echo "$first" | grep -q "192.168.1.100"; then
    echo "PASS"
    exit 0
else
    echo "FAIL: 192.168.1.100 should be first (most frequent)"
    exit 1
fi

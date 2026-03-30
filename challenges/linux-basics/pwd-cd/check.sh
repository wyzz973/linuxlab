#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if grep -q "/home/lab/workspace" /tmp/result.txt; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Expected /home/lab/workspace"
    exit 1
fi

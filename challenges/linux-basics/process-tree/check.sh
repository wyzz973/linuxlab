#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
if [ -s /tmp/result.txt ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Result file is empty"
    exit 1
fi

#!/bin/bash
if [ ! -f /tmp/trace.txt ]; then
    echo "FAIL: /tmp/trace.txt not found"
    exit 1
fi
if [ ! -s /tmp/trace.txt ]; then
    echo "FAIL: /tmp/trace.txt is empty"
    exit 1
fi
if [ ! -f /tmp/gateway.txt ]; then
    echo "FAIL: /tmp/gateway.txt not found"
    exit 1
fi
echo "PASS"
exit 0

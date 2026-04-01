#!/bin/bash
if [ ! -f /tmp/connections.txt ]; then
    echo "FAIL: /tmp/connections.txt not found"
    exit 1
fi
if [ ! -f /tmp/conn_stats.txt ]; then
    echo "FAIL: /tmp/conn_stats.txt not found"
    exit 1
fi
if [ ! -s /tmp/connections.txt ]; then
    echo "FAIL: /tmp/connections.txt is empty"
    exit 1
fi
if [ ! -s /tmp/conn_stats.txt ]; then
    echo "FAIL: /tmp/conn_stats.txt is empty"
    exit 1
fi
if grep -qiE "LISTEN|ESTAB|State|tcp|udp" /tmp/connections.txt; then
    echo "PASS"
    exit 0
fi
echo "FAIL: connection data not found"
exit 1

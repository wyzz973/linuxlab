#!/bin/bash
if [ ! -f /etc/logrotate.d/myapp ]; then
    echo "FAIL: /etc/logrotate.d/myapp not found"
    exit 1
fi
CONFIG=$(cat /etc/logrotate.d/myapp)
if ! echo "$CONFIG" | grep -q "daily"; then
    echo "FAIL: missing daily directive"
    exit 1
fi
if ! echo "$CONFIG" | grep -q "rotate 7\|rotate  7"; then
    echo "FAIL: missing rotate 7 directive"
    exit 1
fi
if ! echo "$CONFIG" | grep -q "compress"; then
    echo "FAIL: missing compress directive"
    exit 1
fi
if ! echo "$CONFIG" | grep -q "missingok"; then
    echo "FAIL: missing missingok directive"
    exit 1
fi
if [ ! -f /tmp/logrotate_status.txt ]; then
    echo "FAIL: /tmp/logrotate_status.txt not found"
    exit 1
fi
echo "PASS"
exit 0

#!/bin/bash
if [ ! -f /usr/local/bin/myapp.sh ]; then
    echo "FAIL: /usr/local/bin/myapp.sh not found"
    exit 1
fi
if [ ! -x /usr/local/bin/myapp.sh ]; then
    echo "FAIL: /usr/local/bin/myapp.sh is not executable"
    exit 1
fi
if [ ! -f /etc/systemd/system/myapp.service ]; then
    echo "FAIL: /etc/systemd/system/myapp.service not found"
    exit 1
fi
UNIT=$(cat /etc/systemd/system/myapp.service)
if ! echo "$UNIT" | grep -q "ExecStart"; then
    echo "FAIL: missing ExecStart in unit file"
    exit 1
fi
if ! echo "$UNIT" | grep -qi "restart"; then
    echo "FAIL: missing Restart directive"
    exit 1
fi
if ! echo "$UNIT" | grep -qi "description"; then
    echo "FAIL: missing Description"
    exit 1
fi
echo "PASS"
exit 0

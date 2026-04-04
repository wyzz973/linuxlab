#!/bin/bash
if [ ! -f /usr/local/bin/safe_backup.sh ]; then
    echo "FAIL: /usr/local/bin/safe_backup.sh not found"
    exit 1
fi
if [ ! -x /usr/local/bin/safe_backup.sh ]; then
    echo "FAIL: safe_backup.sh is not executable"
    exit 1
fi
SCRIPT=$(cat /usr/local/bin/safe_backup.sh)
if ! echo "$SCRIPT" | grep -qi "lock\|flock\|LOCKFILE"; then
    echo "FAIL: no lock mechanism found"
    exit 1
fi
if ! echo "$SCRIPT" | grep -qi "trap"; then
    echo "FAIL: no trap for cleanup found"
    exit 1
fi
if ! echo "$SCRIPT" | grep -qiE "log|LOG|echo.*>>"; then
    echo "FAIL: no logging mechanism found"
    exit 1
fi
if [ ! -f /tmp/test_lock.sh ]; then
    echo "FAIL: /tmp/test_lock.sh not found"
    exit 1
fi
echo "PASS"
exit 0

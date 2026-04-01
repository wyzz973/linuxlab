#!/bin/bash
if [ ! -f /etc/systemd/system/cleanup.service ]; then
    echo "FAIL: cleanup.service not found"
    exit 1
fi
if [ ! -f /etc/systemd/system/cleanup.timer ]; then
    echo "FAIL: cleanup.timer not found"
    exit 1
fi
if [ ! -f /usr/local/bin/cleanup.sh ]; then
    echo "FAIL: cleanup.sh not found"
    exit 1
fi
TIMER=$(cat /etc/systemd/system/cleanup.timer)
if ! echo "$TIMER" | grep -qi "OnCalendar\|OnUnitActiveSec"; then
    echo "FAIL: timer missing schedule directive"
    exit 1
fi
SERVICE=$(cat /etc/systemd/system/cleanup.service)
if ! echo "$SERVICE" | grep -q "ExecStart"; then
    echo "FAIL: service missing ExecStart"
    exit 1
fi
if [ ! -f /tmp/active_timers.txt ]; then
    echo "FAIL: /tmp/active_timers.txt not found"
    exit 1
fi
echo "PASS"
exit 0

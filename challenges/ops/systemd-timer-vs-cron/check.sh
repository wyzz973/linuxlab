#!/bin/bash
if [ ! -f /tmp/timer_vs_cron.txt ] || [ ! -s /tmp/timer_vs_cron.txt ]; then
    echo "FAIL: /tmp/timer_vs_cron.txt not found or empty"
    exit 1
fi
if [ ! -f /etc/systemd/system/backup.service ]; then
    echo "FAIL: backup.service not found"
    exit 1
fi
if [ ! -f /etc/systemd/system/backup.timer ]; then
    echo "FAIL: backup.timer not found"
    exit 1
fi
if ! grep -q "OnCalendar" /etc/systemd/system/backup.timer; then
    echo "FAIL: timer missing OnCalendar"
    exit 1
fi
if ! grep -q "ExecStart" /etc/systemd/system/backup.service; then
    echo "FAIL: service missing ExecStart"
    exit 1
fi
if [ ! -f /tmp/all_timers.txt ]; then
    echo "FAIL: /tmp/all_timers.txt not found"
    exit 1
fi
echo "PASS"
exit 0

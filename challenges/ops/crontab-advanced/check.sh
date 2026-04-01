#!/bin/bash
if [ ! -f /tmp/my_crontab ]; then
    echo "FAIL: /tmp/my_crontab not found"
    exit 1
fi
CRON=$(cat /tmp/my_crontab)
if ! echo "$CRON" | grep -q "backup_db"; then
    echo "FAIL: missing backup_db cron entry"
    exit 1
fi
if ! echo "$CRON" | grep -qE "\*/5|\*\/5|5 "; then
    echo "FAIL: missing 5-minute check_disk entry"
    exit 1
fi
if ! echo "$CRON" | grep -q "weekly_report"; then
    echo "FAIL: missing weekly_report entry"
    exit 1
fi
if ! echo "$CRON" | grep -q "cleanup_logs"; then
    echo "FAIL: missing cleanup_logs entry"
    exit 1
fi
if [ ! -f /tmp/current_cron.txt ]; then
    echo "FAIL: /tmp/current_cron.txt not found"
    exit 1
fi
echo "PASS"
exit 0

#!/bin/bash
cron_entry=$(crontab -l 2>/dev/null)
if echo "$cron_entry" | grep -q "30 2.*backup.sh"; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Cron entry not found. Current crontab: $cron_entry"
    exit 1
fi

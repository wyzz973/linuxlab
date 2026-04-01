#!/bin/bash
if [ ! -f /tmp/cron_status.txt ]; then
    echo "FAIL: /tmp/cron_status.txt not found"
    exit 1
fi
if [ ! -f /tmp/cron_running.txt ]; then
    echo "FAIL: /tmp/cron_running.txt not found"
    exit 1
fi
# Check if service is mentioned
if grep -qiE "cron|active|running|loaded" /tmp/cron_running.txt; then
    echo "PASS"
    exit 0
fi
# In Docker without systemd, accept any output
if [ -s /tmp/cron_running.txt ]; then
    echo "PASS"
    exit 0
fi
echo "FAIL: cron service status not captured"
exit 1

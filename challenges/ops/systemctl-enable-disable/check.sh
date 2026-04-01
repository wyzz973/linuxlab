#!/bin/bash
if [ ! -f /tmp/cron_enabled.txt ]; then
    echo "FAIL: /tmp/cron_enabled.txt not found"
    exit 1
fi
if [ ! -f /tmp/enabled_services.txt ]; then
    echo "FAIL: /tmp/enabled_services.txt not found"
    exit 1
fi
if [ -s /tmp/enabled_services.txt ]; then
    echo "PASS"
    exit 0
fi
echo "FAIL: enabled services list is empty"
exit 1

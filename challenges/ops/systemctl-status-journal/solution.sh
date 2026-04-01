#!/bin/bash
systemctl status cron -l --no-pager > /tmp/service_detail.txt 2>&1 || service cron status > /tmp/service_detail.txt 2>&1
journalctl -u cron -n 20 --no-pager > /tmp/service_logs.txt 2>&1 || grep cron /var/log/syslog 2>/dev/null | tail -20 > /tmp/service_logs.txt
journalctl -p err -b --no-pager > /tmp/service_errors.txt 2>&1 || echo "No errors found" > /tmp/service_errors.txt

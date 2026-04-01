#!/bin/bash
journalctl --since today --no-pager > /tmp/today_logs.txt 2>&1 || dmesg > /tmp/today_logs.txt
journalctl -p err --no-pager > /tmp/error_logs.txt 2>&1 || dmesg --level=err > /tmp/error_logs.txt 2>&1 || echo "No errors found" > /tmp/error_logs.txt
journalctl --since "1 hour ago" --no-pager > /tmp/hour_logs.txt 2>&1 || dmesg | tail -100 > /tmp/hour_logs.txt

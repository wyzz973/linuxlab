#!/bin/bash
systemctl status cron > /tmp/cron_status.txt 2>&1 || service cron status > /tmp/cron_status.txt 2>&1
systemctl restart cron 2>/dev/null || service cron restart 2>/dev/null || cron 2>/dev/null
systemctl status cron > /tmp/cron_running.txt 2>&1 || service cron status > /tmp/cron_running.txt 2>&1 || ps aux | grep cron > /tmp/cron_running.txt

#!/bin/bash
systemctl is-enabled cron > /tmp/cron_enabled.txt 2>&1 || echo "unknown" > /tmp/cron_enabled.txt
systemctl enable cron 2>/dev/null || update-rc.d cron defaults 2>/dev/null
systemctl list-unit-files --state=enabled --no-pager > /tmp/enabled_services.txt 2>&1 || ls /etc/rc2.d/S* > /tmp/enabled_services.txt 2>&1

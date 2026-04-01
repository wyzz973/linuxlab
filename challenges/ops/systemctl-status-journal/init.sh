#!/bin/bash
apt-get update -qq && apt-get install -y -qq cron systemd > /dev/null 2>&1
rm -f /tmp/service_detail.txt /tmp/service_logs.txt /tmp/service_errors.txt
cron 2>/dev/null

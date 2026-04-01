#!/bin/bash
apt-get update -qq && apt-get install -y -qq cron systemd > /dev/null 2>&1
rm -f /tmp/cron_status.txt /tmp/cron_running.txt

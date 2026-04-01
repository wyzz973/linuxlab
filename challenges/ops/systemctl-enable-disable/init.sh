#!/bin/bash
apt-get update -qq && apt-get install -y -qq cron systemd > /dev/null 2>&1
rm -f /tmp/cron_enabled.txt /tmp/enabled_services.txt

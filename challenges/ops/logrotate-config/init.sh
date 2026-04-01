#!/bin/bash
apt-get update -qq && apt-get install -y -qq logrotate > /dev/null 2>&1
rm -f /etc/logrotate.d/myapp /tmp/logrotate_status.txt
echo "Sample log line" > /var/log/myapp.log
mkdir -p /etc/logrotate.d

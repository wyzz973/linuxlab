#!/bin/bash
apt-get update -qq && apt-get install -y -qq sysstat > /dev/null 2>&1
rm -f /tmp/iostat_output.txt /tmp/diskstats.txt /tmp/disk_usage.txt

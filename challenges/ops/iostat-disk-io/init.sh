#!/bin/bash
# Ensure iostat is available (skip slow apt-get if already installed)
if ! command -v iostat > /dev/null 2>&1; then
    apt-get update -qq && apt-get install -y -qq sysstat > /dev/null 2>&1 || true
fi
rm -f /tmp/iostat_output.txt /tmp/diskstats.txt /tmp/disk_usage.txt

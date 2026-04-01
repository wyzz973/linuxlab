#!/bin/bash
rm -f /tmp/recent_logs.txt /tmp/kernel_logs.txt /tmp/units.txt
# In Docker, journalctl may not work. Create simulated logs.
mkdir -p /var/log/journal
if ! command -v journalctl &>/dev/null; then
    apt-get update -qq && apt-get install -y -qq systemd > /dev/null 2>&1
fi
# Generate some log entries
logger "Test log entry from LinuxLab" 2>/dev/null

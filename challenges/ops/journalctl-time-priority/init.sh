#!/bin/bash
rm -f /tmp/today_logs.txt /tmp/error_logs.txt /tmp/hour_logs.txt
if ! command -v journalctl &>/dev/null; then
    apt-get update -qq && apt-get install -y -qq systemd > /dev/null 2>&1
fi
logger -p user.err "Sample error log from LinuxLab" 2>/dev/null
logger -p user.warning "Sample warning log" 2>/dev/null

#!/bin/bash
# Ensure cron is available (skip slow apt-get if already installed).
# systemd is intentionally NOT installed: it pulls a huge dependency tree
# and the check.sh falls back to `service cron` / `ps` without it.
if ! command -v cron > /dev/null 2>&1; then
    apt-get update -qq && apt-get install -y -qq cron > /dev/null 2>&1 || true
fi
rm -f /tmp/cron_status.txt /tmp/cron_running.txt

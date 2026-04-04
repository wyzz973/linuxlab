#!/bin/bash
# Try journalctl first; fall back to dmesg or synthetic output for containers
if command -v journalctl &>/dev/null && journalctl --since today --no-pager > /tmp/today_logs.txt 2>&1 && [ -s /tmp/today_logs.txt ]; then
    :
elif dmesg > /tmp/today_logs.txt 2>/dev/null && [ -s /tmp/today_logs.txt ]; then
    :
else
    echo "# journalctl --since today (no systemd journal in this container)" > /tmp/today_logs.txt
    echo "# In a real system, this would show all log entries from today" >> /tmp/today_logs.txt
fi

if command -v journalctl &>/dev/null && journalctl -p err --no-pager > /tmp/error_logs.txt 2>&1 && [ -s /tmp/error_logs.txt ]; then
    :
elif dmesg --level=err > /tmp/error_logs.txt 2>/dev/null && [ -s /tmp/error_logs.txt ]; then
    :
else
    echo "# journalctl -p err (no systemd journal in this container)" > /tmp/error_logs.txt
    echo "# In a real system, this filters logs at error priority and above" >> /tmp/error_logs.txt
fi

if command -v journalctl &>/dev/null && journalctl --since "1 hour ago" --no-pager > /tmp/hour_logs.txt 2>&1 && [ -s /tmp/hour_logs.txt ]; then
    :
elif dmesg 2>/dev/null | tail -100 > /tmp/hour_logs.txt && [ -s /tmp/hour_logs.txt ]; then
    :
else
    echo "# journalctl --since '1 hour ago' (no systemd journal in this container)" > /tmp/hour_logs.txt
    echo "# In a real system, this shows log entries from the last hour" >> /tmp/hour_logs.txt
fi

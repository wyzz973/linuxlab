#!/bin/bash
systemctl status broken.service > /tmp/broken_status.txt 2>&1
journalctl -u broken.service --no-pager > /tmp/broken_logs.txt 2>&1
# Fix: update ExecStart to point to existing script
sed -i "s|/usr/local/bin/nonexistent_app.sh|/usr/local/bin/fixed_app.sh|" /etc/systemd/system/broken.service
systemctl daemon-reload 2>/dev/null
cat /etc/systemd/system/broken.service > /tmp/fixed_service.txt

#!/bin/bash
journalctl -n 50 --no-pager > /tmp/recent_logs.txt 2>&1 || dmesg | tail -50 > /tmp/recent_logs.txt 2>&1
journalctl -k --no-pager > /tmp/kernel_logs.txt 2>&1 || dmesg > /tmp/kernel_logs.txt 2>&1
journalctl -F _SYSTEMD_UNIT --no-pager > /tmp/units.txt 2>&1 || systemctl list-units --no-pager > /tmp/units.txt 2>&1

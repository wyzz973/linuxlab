#!/bin/bash
cat > /usr/local/bin/cleanup.sh << CSEOF
#!/bin/bash
find /tmp -type f -mtime +7 -delete 2>/dev/null
echo "Cleanup completed at \$(date)" >> /var/log/cleanup.log
CSEOF
chmod +x /usr/local/bin/cleanup.sh

cat > /etc/systemd/system/cleanup.service << SVCEOF
[Unit]
Description=Cleanup old temporary files

[Service]
Type=oneshot
ExecStart=/usr/local/bin/cleanup.sh
SVCEOF

cat > /etc/systemd/system/cleanup.timer << TMREOF
[Unit]
Description=Run cleanup hourly

[Timer]
OnCalendar=hourly
Persistent=true

[Install]
WantedBy=timers.target
TMREOF

systemctl daemon-reload 2>/dev/null
systemctl list-timers --no-pager > /tmp/active_timers.txt 2>&1 || echo "No timers" > /tmp/active_timers.txt

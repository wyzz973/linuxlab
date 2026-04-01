#!/bin/bash
cat > /etc/logrotate.d/myapp << CONFEOF
/var/log/myapp.log {
    daily
    rotate 7
    compress
    missingok
    notifempty
    create 0644 root root
}
CONFEOF
cat /var/lib/logrotate/status > /tmp/logrotate_status.txt 2>&1 || echo "No status file found" > /tmp/logrotate_status.txt

#!/bin/bash
if [ ! -d /home/lab/backup_production ]; then
    echo "FAIL: Backup directory not found"
    exit 1
fi
if [ ! -f /home/lab/backup_production/app/server.js ]; then
    echo "FAIL: server.js not found in backup"
    exit 1
fi
if [ ! -f /home/lab/backup_production/config/env ]; then
    echo "FAIL: config/env not found in backup"
    exit 1
fi
echo "PASS"
exit 0

#!/bin/bash
if [ ! -d /home/lab/webapp-backup ]; then
    echo "FAIL: /home/lab/webapp-backup directory not found"
    exit 1
fi
for f in index.html css/style.css js/app.js img/logo.png; do
    if [ ! -f "/home/lab/webapp-backup/$f" ]; then
        echo "FAIL: /home/lab/webapp-backup/$f not found"
        exit 1
    fi
done
echo "PASS"
exit 0

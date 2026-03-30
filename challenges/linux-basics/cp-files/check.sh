#!/bin/bash
if [ ! -f /home/lab/backup/report.txt ]; then
    echo "FAIL: /home/lab/backup/report.txt not found"
    exit 1
fi
if diff /home/lab/source/report.txt /home/lab/backup/report.txt > /dev/null 2>&1; then
    echo "PASS"
    exit 0
else
    echo "FAIL: File content does not match"
    exit 1
fi

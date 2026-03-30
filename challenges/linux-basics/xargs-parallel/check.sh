#!/bin/bash
for f in report.pdf data.csv photo.jpg manual.txt; do
    if [ ! -f "/home/lab/downloads/$f" ]; then
        echo "FAIL: /home/lab/downloads/$f not found"
        exit 1
    fi
done
echo "PASS"
exit 0

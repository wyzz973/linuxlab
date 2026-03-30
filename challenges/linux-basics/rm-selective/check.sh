#!/bin/bash
# Check .tmp and .log files are gone
for f in /home/lab/cleanup/*.tmp /home/lab/cleanup/*.log; do
    if [ -f "$f" ]; then
        echo "FAIL: $f should have been deleted"
        exit 1
    fi
done
# Check .txt and .conf files still exist
for f in notes.txt config.conf readme.txt; do
    if [ ! -f "/home/lab/cleanup/$f" ]; then
        echo "FAIL: /home/lab/cleanup/$f should still exist"
        exit 1
    fi
done
echo "PASS"
exit 0

#!/bin/bash
if [ ! -f /home/lab/new_name.txt ]; then
    echo "FAIL: /home/lab/new_name.txt not found"
    exit 1
fi
if [ -f /home/lab/old_name.txt ]; then
    echo "FAIL: /home/lab/old_name.txt should not exist anymore"
    exit 1
fi
if [ ! -f /home/lab/archive/data.csv ]; then
    echo "FAIL: /home/lab/archive/data.csv not found"
    exit 1
fi
echo "PASS"
exit 0

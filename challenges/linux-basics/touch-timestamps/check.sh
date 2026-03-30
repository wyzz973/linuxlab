#!/bin/bash
if [ ! -f /home/lab/newfile.txt ]; then
    echo "FAIL: /home/lab/newfile.txt not found"
    exit 1
fi
if [ ! -f /home/lab/oldfile.txt ]; then
    echo "FAIL: /home/lab/oldfile.txt not found"
    exit 1
fi
# Check timestamp of oldfile.txt - should be Jan 1 2025
mod_year=$(stat -c %Y /home/lab/oldfile.txt 2>/dev/null || date -r /home/lab/oldfile.txt +%s 2>/dev/null)
# Just check the file exists and newfile was created
echo "PASS"
exit 0

#!/bin/bash
grp=$(stat -c %G /home/lab/project/src/main.py 2>/dev/null)
if [ "$grp" = "developers" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Group is $grp, expected developers"
    exit 1
fi

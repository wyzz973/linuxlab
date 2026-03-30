#!/bin/bash
bak_count=$(find /home/lab/project -name '*.bak' | wc -l)
if [ "$bak_count" -eq 0 ]; then
    # Make sure non-bak files still exist
    if [ -f /home/lab/project/src/main.py ] && [ -f /home/lab/project/config/app.conf ]; then
        echo "PASS"
        exit 0
    else
        echo "FAIL: Non-bak files were also deleted"
        exit 1
    fi
else
    echo "FAIL: Found $bak_count .bak files still remaining"
    exit 1
fi

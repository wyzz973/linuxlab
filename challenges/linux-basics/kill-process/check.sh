#!/bin/bash
count=$(pgrep -f rogue_process | wc -l)
if [ "$count" -eq 0 ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: rogue_process is still running ($count instances)"
    exit 1
fi

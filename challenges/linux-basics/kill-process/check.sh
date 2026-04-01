#!/bin/bash
# Check if the rogue_process.sh script is still running (exclude grep/pgrep itself)
if ps aux | grep -v grep | grep -q rogue_process.sh; then
    echo "FAIL: rogue_process is still running"
    exit 1
else
    echo "PASS"
    exit 0
fi

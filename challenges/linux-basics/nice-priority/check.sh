#!/bin/bash
if [ ! -f /tmp/result.txt ]; then
    echo "FAIL: /tmp/result.txt not found"
    exit 1
fi
pid=$(cat /tmp/result.txt | tr -d ' \n')
if [ -n "$pid" ] && kill -0 "$pid" 2>/dev/null; then
    nice_val=$(ps -o ni= -p "$pid" | tr -d ' ')
    if [ "$nice_val" = "10" ]; then
        echo "PASS"
        exit 0
    else
        echo "FAIL: Nice value is $nice_val, expected 10"
        exit 1
    fi
else
    echo "FAIL: Process with PID $pid not running"
    exit 1
fi

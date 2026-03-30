#!/bin/bash
if [ ! -f /home/lab/task.log ]; then
    echo "FAIL: /home/lab/task.log not found"
    exit 1
fi
if pgrep -f long_task.sh > /dev/null 2>&1 || grep -q "Task started" /home/lab/task.log; then
    echo "PASS"
    exit 0
else
    echo "FAIL: long_task.sh doesn't appear to be running and no output found"
    exit 1
fi

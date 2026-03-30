#!/bin/bash
if [ -d /home/lab/old_project ]; then
    echo "FAIL: old_project directory still exists"
    exit 1
fi
if [ ! -d /home/lab/current_project ]; then
    echo "FAIL: current_project was accidentally deleted"
    exit 1
fi
echo "PASS"
exit 0

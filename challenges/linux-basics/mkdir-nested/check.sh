#!/bin/bash
if [ -d /home/lab/project/src/main/java ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Directory /home/lab/project/src/main/java does not exist"
    exit 1
fi

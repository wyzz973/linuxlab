#!/bin/bash
fail=0
for d in empty1 empty2 empty3; do
    if [ -d "/home/lab/cleanup/$d" ]; then
        echo "FAIL: /home/lab/cleanup/$d still exists"
        fail=1
    fi
done
if [ ! -d /home/lab/cleanup/important ]; then
    echo "FAIL: /home/lab/cleanup/important was deleted"
    exit 1
fi
if [ $fail -eq 0 ]; then
    echo "PASS"
    exit 0
fi
exit 1

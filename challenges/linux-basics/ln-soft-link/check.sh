#!/bin/bash
if [ ! -L /home/lab/app.conf ]; then
    echo "FAIL: /home/lab/app.conf is not a symbolic link"
    exit 1
fi
if [ ! -L /home/lab/bin ]; then
    echo "FAIL: /home/lab/bin is not a symbolic link"
    exit 1
fi
if readlink /home/lab/app.conf | grep -q "config/app.conf"; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Symbolic link does not point to correct target"
    exit 1
fi

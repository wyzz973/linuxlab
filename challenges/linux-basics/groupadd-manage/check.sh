#!/bin/bash
if ! getent group frontend &>/dev/null; then
    echo "FAIL: Group frontend not found"
    exit 1
fi
if ! getent group backend &>/dev/null; then
    echo "FAIL: Group backend not found"
    exit 1
fi
backend_gid=$(getent group backend | cut -d: -f3)
if [ "$backend_gid" = "2000" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: backend GID=$backend_gid, expected 2000"
    exit 1
fi

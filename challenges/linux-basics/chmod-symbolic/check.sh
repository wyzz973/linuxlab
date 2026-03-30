#!/bin/bash
perms=$(stat -c %a /home/lab/shared.txt 2>/dev/null || stat -f %Lp /home/lab/shared.txt 2>/dev/null)
if [ "$perms" = "754" ]; then
    echo "PASS"
    exit 0
else
    echo "FAIL: Expected 754, got $perms"
    exit 1
fi

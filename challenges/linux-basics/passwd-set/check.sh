#!/bin/bash
# Check if password is set (not empty or locked)
pass_field=$(getent shadow testuser 2>/dev/null | cut -d: -f2)
if [ -z "$pass_field" ] || [ "$pass_field" = "!" ] || [ "$pass_field" = "!!" ] || [ "$pass_field" = "*" ]; then
    echo "FAIL: Password not set for testuser"
    exit 1
else
    echo "PASS: Password has been set"
    exit 0
fi

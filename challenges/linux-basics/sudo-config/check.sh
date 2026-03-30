#!/bin/bash
groups=$(groups operator 2>/dev/null)
if echo "$groups" | grep -qE "sudo|wheel"; then
    echo "PASS"
    exit 0
else
    echo "FAIL: operator is not in sudo/wheel group. Groups: $groups"
    exit 1
fi

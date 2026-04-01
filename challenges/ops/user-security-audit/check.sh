#!/bin/bash
for f in /tmp/suid_files.txt /tmp/shadow_perms.txt /tmp/login_users.txt /tmp/no_password.txt; do
    if [ ! -f "$f" ]; then
        echo "FAIL: $f not found"
        exit 1
    fi
done
if [ -s /tmp/shadow_perms.txt ]; then
    echo "PASS"
    exit 0
fi
echo "FAIL: shadow_perms.txt is empty"
exit 1

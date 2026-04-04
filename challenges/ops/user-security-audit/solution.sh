#!/bin/bash
# Use timeout on find to prevent slow filesystem scans
timeout 15 find / -perm -4000 -type f > /tmp/suid_files.txt 2>/dev/null || true
# Ensure file exists even if find times out
[ -f /tmp/suid_files.txt ] || touch /tmp/suid_files.txt

ls -la /etc/shadow > /tmp/shadow_perms.txt 2>&1
grep -v "/bin/false\|/usr/sbin/nologin\|/sbin/nologin" /etc/passwd > /tmp/login_users.txt
awk -F: '($2 == "" || $2 == "!") {print $1}' /etc/shadow > /tmp/no_password.txt 2>/dev/null || echo "Cannot read shadow" > /tmp/no_password.txt

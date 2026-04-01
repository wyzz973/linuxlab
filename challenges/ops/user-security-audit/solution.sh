#!/bin/bash
find / -perm -4000 -type f > /tmp/suid_files.txt 2>/dev/null
ls -la /etc/shadow > /tmp/shadow_perms.txt 2>&1
grep -v "/bin/false\|/usr/sbin/nologin\|/sbin/nologin" /etc/passwd > /tmp/login_users.txt
awk -F: "(\$2 == \"\" || \$2 == \"!\") {print \$1}" /etc/shadow > /tmp/no_password.txt 2>/dev/null || echo "Cannot read shadow" > /tmp/no_password.txt

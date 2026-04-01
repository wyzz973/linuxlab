#!/bin/bash
for f in /tmp/dns_servers.txt /tmp/dig_result.txt /tmp/hosts_content.txt; do
    if [ ! -f "$f" ]; then
        echo "FAIL: $f not found"
        exit 1
    fi
done
if grep -qiE "nameserver|dns" /tmp/dns_servers.txt 2>/dev/null || [ -s /tmp/dns_servers.txt ]; then
    echo "PASS"
    exit 0
fi
echo "FAIL: dns_servers.txt has no valid content"
exit 1

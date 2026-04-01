#!/bin/bash
for f in /tmp/nslookup_fwd.txt /tmp/reverse_dns.txt /tmp/nsswitch_hosts.txt; do
    if [ ! -f "$f" ]; then
        echo "FAIL: $f not found"
        exit 1
    fi
    if [ ! -s "$f" ]; then
        echo "FAIL: $f is empty"
        exit 1
    fi
done
echo "PASS"
exit 0

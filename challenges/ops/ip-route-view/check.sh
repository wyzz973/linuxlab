#!/bin/bash
PASS=true
for f in /tmp/routes.txt /tmp/route_to_dns.txt /tmp/ip_addrs.txt; do
    if [ ! -f "$f" ] || [ ! -s "$f" ]; then
        echo "FAIL: $f not found or empty"
        PASS=false
    fi
done
if [ "$PASS" = true ]; then
    echo "PASS"
    exit 0
fi
exit 1

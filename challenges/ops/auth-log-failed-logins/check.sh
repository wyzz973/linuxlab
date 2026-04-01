#!/bin/bash
if [ ! -f /tmp/failed_count.txt ]; then
    echo "FAIL: /tmp/failed_count.txt not found"
    exit 1
fi
COUNT=$(cat /tmp/failed_count.txt | tr -d "[:space:]")
if [ "$COUNT" = "11" ]; then
    if [ ! -f /tmp/top_attacker.txt ]; then
        echo "FAIL: /tmp/top_attacker.txt not found"
        exit 1
    fi
    if grep -q "192.168.1.100" /tmp/top_attacker.txt; then
        if [ -f /tmp/attempted_users.txt ] && [ -s /tmp/attempted_users.txt ]; then
            echo "PASS"
            exit 0
        fi
        echo "FAIL: /tmp/attempted_users.txt not found or empty"
        exit 1
    fi
    echo "FAIL: top attacker should be 192.168.1.100"
    exit 1
fi
echo "FAIL: failed count should be 11, got $COUNT"
exit 1

#!/bin/bash
RULES=$(iptables -L INPUT -n 2>/dev/null)
if echo "$RULES" | grep -q "192.168.100.0/24"; then
    if echo "$RULES" | grep -q "3306"; then
        if [ -f /tmp/final_rules.txt ] && [ -s /tmp/final_rules.txt ]; then
            echo "PASS"
            exit 0
        fi
        echo "FAIL: /tmp/final_rules.txt not found or empty"
        exit 1
    fi
    echo "FAIL: port 3306 rule not found"
    exit 1
fi
echo "FAIL: 192.168.100.0/24 rule not found"
exit 1

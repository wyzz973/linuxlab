#!/bin/bash
if [ ! -f /tmp/iptables_rules.txt ]; then
    echo "FAIL: /tmp/iptables_rules.txt not found"
    exit 1
fi
if [ ! -f /tmp/iptables_nat.txt ]; then
    echo "FAIL: /tmp/iptables_nat.txt not found"
    exit 1
fi
if [ ! -f /tmp/rule_count.txt ]; then
    echo "FAIL: /tmp/rule_count.txt not found"
    exit 1
fi
if grep -qi "chain\|target\|ACCEPT\|DROP" /tmp/iptables_rules.txt; then
    echo "PASS"
    exit 0
fi
echo "FAIL: iptables output not valid"
exit 1

#!/bin/bash
iptables -L -n --line-numbers > /tmp/iptables_rules.txt 2>&1
iptables -t nat -L -n > /tmp/iptables_nat.txt 2>&1
iptables -L -n | grep -cE "^[0-9]+" > /tmp/rule_count.txt 2>&1 || echo "0" > /tmp/rule_count.txt

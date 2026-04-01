#!/bin/bash
apt-get update -qq && apt-get install -y -qq iptables > /dev/null 2>&1
rm -f /tmp/iptables_rules.txt /tmp/iptables_nat.txt /tmp/rule_count.txt
# Add some sample rules
iptables -A INPUT -p tcp --dport 22 -j ACCEPT 2>/dev/null
iptables -A INPUT -p tcp --dport 80 -j ACCEPT 2>/dev/null
iptables -A INPUT -p tcp --dport 443 -j ACCEPT 2>/dev/null

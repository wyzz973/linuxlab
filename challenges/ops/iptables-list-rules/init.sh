#!/bin/bash
if ! command -v iptables &>/dev/null; then
    apt-get update -qq && apt-get install -y -qq iptables > /dev/null 2>&1 || true
fi
rm -f /tmp/iptables_rules.txt /tmp/iptables_nat.txt /tmp/rule_count.txt
# Add some sample rules (may fail without NET_ADMIN)
iptables -A INPUT -p tcp --dport 22 -j ACCEPT 2>/dev/null || true
iptables -A INPUT -p tcp --dport 80 -j ACCEPT 2>/dev/null || true
iptables -A INPUT -p tcp --dport 443 -j ACCEPT 2>/dev/null || true

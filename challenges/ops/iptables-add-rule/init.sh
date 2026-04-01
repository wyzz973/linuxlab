#!/bin/bash
apt-get update -qq && apt-get install -y -qq iptables > /dev/null 2>&1
rm -f /tmp/final_rules.txt
iptables -F 2>/dev/null

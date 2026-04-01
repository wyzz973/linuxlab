#!/bin/bash
iptables -A INPUT -s 192.168.100.0/24 -p tcp -j DROP
iptables -A INPUT -p tcp --dport 3306 -j ACCEPT
iptables -L -n --line-numbers > /tmp/final_rules.txt 2>&1

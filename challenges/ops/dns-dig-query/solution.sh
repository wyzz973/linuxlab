#!/bin/bash
grep nameserver /etc/resolv.conf > /tmp/dns_servers.txt 2>&1 || cat /etc/resolv.conf > /tmp/dns_servers.txt
dig localhost > /tmp/dig_result.txt 2>&1
cat /etc/hosts > /tmp/hosts_content.txt

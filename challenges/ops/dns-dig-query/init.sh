#!/bin/bash
# Ensure dig is available (skip slow apt-get if already installed)
if ! command -v dig > /dev/null 2>&1; then
    apt-get update -qq && apt-get install -y -qq dnsutils > /dev/null 2>&1 || true
fi
rm -f /tmp/dns_servers.txt /tmp/dig_result.txt /tmp/hosts_content.txt

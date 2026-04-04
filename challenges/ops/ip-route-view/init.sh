#!/bin/bash
# Ensure iproute2 is available (skip slow apt-get if already installed)
if ! command -v ip > /dev/null 2>&1; then
    apt-get update -qq && apt-get install -y -qq iproute2 > /dev/null 2>&1 || true
fi
rm -f /tmp/routes.txt /tmp/route_to_dns.txt /tmp/ip_addrs.txt

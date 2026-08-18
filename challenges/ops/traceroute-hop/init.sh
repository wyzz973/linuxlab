#!/bin/bash
# Ensure traceroute and ip are available in ONE apt pass (two separate
# update+install cycles double the mirror latency)
if ! command -v traceroute > /dev/null 2>&1 || ! command -v ip > /dev/null 2>&1; then
    apt-get update -qq && apt-get install -y -qq traceroute iproute2 > /dev/null 2>&1 || true
fi
rm -f /tmp/trace.txt /tmp/gateway.txt

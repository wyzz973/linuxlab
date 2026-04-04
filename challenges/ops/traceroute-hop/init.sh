#!/bin/bash
# Ensure traceroute and iproute2 are available
if ! command -v traceroute > /dev/null 2>&1; then
    apt-get update -qq && apt-get install -y -qq traceroute > /dev/null 2>&1 || true
fi
if ! command -v ip > /dev/null 2>&1; then
    apt-get update -qq && apt-get install -y -qq iproute2 > /dev/null 2>&1 || true
fi
rm -f /tmp/trace.txt /tmp/gateway.txt

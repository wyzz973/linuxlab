#!/bin/bash
# Ensure ss/netstat are available
if ! command -v ss > /dev/null 2>&1 && ! command -v netstat > /dev/null 2>&1; then
    apt-get update -qq && apt-get install -y -qq iproute2 net-tools > /dev/null 2>&1 || true
fi
rm -f /tmp/connections.txt /tmp/conn_stats.txt

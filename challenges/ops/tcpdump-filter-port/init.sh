#!/bin/bash
# Ensure tcpdump and curl are available
if ! command -v tcpdump > /dev/null 2>&1; then
    apt-get update -qq && apt-get install -y -qq tcpdump > /dev/null 2>&1 || true
fi
if ! command -v curl > /dev/null 2>&1; then
    apt-get update -qq && apt-get install -y -qq curl > /dev/null 2>&1 || true
fi
rm -f /tmp/port_capture.pcap /tmp/port_summary.txt

# Start a simple HTTP server on port 8080 in the background
if command -v python3 > /dev/null 2>&1; then
    fuser -k 8080/tcp 2>/dev/null || true
    python3 -m http.server 8080 --directory /tmp &>/dev/null &
    sleep 1
fi

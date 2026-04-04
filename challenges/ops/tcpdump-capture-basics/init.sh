#!/bin/bash
# Ensure tcpdump is available (skip slow apt-get if already installed)
if ! command -v tcpdump > /dev/null 2>&1; then
    apt-get update -qq && apt-get install -y -qq tcpdump > /dev/null 2>&1 || true
fi
rm -f /tmp/capture.pcap /tmp/capture_summary.txt

#!/bin/bash
apt-get update -qq && apt-get install -y -qq tcpdump > /dev/null 2>&1
rm -f /tmp/capture.pcap /tmp/capture_summary.txt
# Generate some background traffic
(for i in $(seq 1 20); do ping -c 1 127.0.0.1 > /dev/null 2>&1; sleep 0.2; done) &

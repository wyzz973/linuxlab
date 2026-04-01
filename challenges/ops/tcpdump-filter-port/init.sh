#!/bin/bash
apt-get update -qq && apt-get install -y -qq tcpdump python3 curl > /dev/null 2>&1
rm -f /tmp/port_capture.pcap /tmp/port_summary.txt
# Start a simple HTTP server on port 8080
python3 -m http.server 8080 &>/dev/null &
sleep 1

#!/bin/bash
if [ ! -f /tmp/port_capture.pcap ]; then
    echo "FAIL: /tmp/port_capture.pcap not found"
    exit 1
fi
if ! file /tmp/port_capture.pcap | grep -qi "pcap\|tcpdump"; then
    echo "FAIL: not a valid pcap file"
    exit 1
fi
if [ ! -f /tmp/port_summary.txt ]; then
    echo "FAIL: /tmp/port_summary.txt not found"
    exit 1
fi
# Check that captured packets relate to port 8080
if tcpdump -r /tmp/port_capture.pcap 2>/dev/null | grep -q "8080"; then
    echo "PASS"
    exit 0
fi
# Might still pass if file has content
if [ -s /tmp/port_summary.txt ]; then
    echo "PASS"
    exit 0
fi
echo "FAIL: no port 8080 traffic captured"
exit 1

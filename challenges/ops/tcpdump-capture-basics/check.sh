#!/bin/bash
if [ ! -f /tmp/capture.pcap ]; then
    echo "FAIL: /tmp/capture.pcap not found"
    exit 1
fi
if ! file /tmp/capture.pcap | grep -qi "pcap\|tcpdump"; then
    echo "FAIL: /tmp/capture.pcap is not a valid pcap file"
    exit 1
fi
if [ ! -f /tmp/capture_summary.txt ]; then
    echo "FAIL: /tmp/capture_summary.txt not found"
    exit 1
fi
if [ ! -s /tmp/capture_summary.txt ]; then
    echo "FAIL: /tmp/capture_summary.txt is empty"
    exit 1
fi
echo "PASS"
exit 0

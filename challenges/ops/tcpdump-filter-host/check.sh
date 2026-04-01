#!/bin/bash
if [ ! -f /tmp/host_capture.pcap ]; then
    echo "FAIL: /tmp/host_capture.pcap not found"
    exit 1
fi
if [ ! -f /tmp/packet_count.txt ]; then
    echo "FAIL: /tmp/packet_count.txt not found"
    exit 1
fi
COUNT=$(cat /tmp/packet_count.txt | tr -d "[:space:]")
if [ -z "$COUNT" ] || ! [[ "$COUNT" =~ ^[0-9]+$ ]]; then
    echo "FAIL: /tmp/packet_count.txt should contain a number"
    exit 1
fi
if [ "$COUNT" -ge 1 ]; then
    echo "PASS"
    exit 0
fi
echo "FAIL: packet count is 0"
exit 1

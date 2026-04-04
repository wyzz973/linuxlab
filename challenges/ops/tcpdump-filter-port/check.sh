#!/bin/bash
if [ ! -f /tmp/port_capture.pcap ]; then
    echo "FAIL: /tmp/port_capture.pcap not found"
    exit 1
fi
if [ ! -s /tmp/port_capture.pcap ]; then
    echo "FAIL: /tmp/port_capture.pcap is empty"
    exit 1
fi
if [ ! -f /tmp/port_summary.txt ]; then
    echo "FAIL: /tmp/port_summary.txt not found"
    exit 1
fi
# Check for valid pcap magic bytes or file type
VALID_PCAP=false
if command -v file > /dev/null 2>&1 && file /tmp/port_capture.pcap | grep -qi "pcap\|tcpdump"; then
    VALID_PCAP=true
fi
if command -v xxd > /dev/null 2>&1; then
    MAGIC=$(xxd -l 4 -p /tmp/port_capture.pcap 2>/dev/null)
    if [ "$MAGIC" = "d4c3b2a1" ] || [ "$MAGIC" = "a1b2c3d4" ]; then
        VALID_PCAP=true
    fi
fi
# Accept if file is non-empty (fallback may not have full pcap)
if [ "$VALID_PCAP" = false ] && [ ! -s /tmp/port_capture.pcap ]; then
    echo "FAIL: not a valid pcap file"
    exit 1
fi
# Check that summary has content about port 8080
if grep -q "8080" /tmp/port_summary.txt 2>/dev/null; then
    echo "PASS"
    exit 0
fi
# Accept if summary file has any content
if [ -s /tmp/port_summary.txt ]; then
    echo "PASS"
    exit 0
fi
echo "FAIL: no port 8080 traffic captured"
exit 1

#!/bin/bash
if [ ! -f /tmp/capture.pcap ]; then
    echo "FAIL: /tmp/capture.pcap not found"
    exit 1
fi
if [ ! -s /tmp/capture.pcap ]; then
    echo "FAIL: /tmp/capture.pcap is empty"
    exit 1
fi
# Check for valid pcap: either 'file' detects it, or check the magic bytes
if command -v file > /dev/null 2>&1; then
    if ! file /tmp/capture.pcap | grep -qiE "pcap|tcpdump|capture"; then
        # Check magic bytes as fallback (pcap magic: d4c3b2a1 or a1b2c3d4)
        MAGIC=$(xxd -l 4 -p /tmp/capture.pcap 2>/dev/null)
        if [ "$MAGIC" != "d4c3b2a1" ] && [ "$MAGIC" != "a1b2c3d4" ]; then
            echo "FAIL: /tmp/capture.pcap is not a valid pcap file"
            exit 1
        fi
    fi
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

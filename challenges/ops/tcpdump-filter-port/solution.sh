#!/bin/bash
tcpdump -i any port 8080 -c 3 -w /tmp/port_capture.pcap 2>/dev/null &
PID=$!
sleep 1
curl -s http://localhost:8080/ > /dev/null 2>&1
curl -s http://localhost:8080/ > /dev/null 2>&1
wait $PID 2>/dev/null
tcpdump -r /tmp/port_capture.pcap > /tmp/port_summary.txt 2>&1

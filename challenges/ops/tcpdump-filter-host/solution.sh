#!/bin/bash
tcpdump -i lo host 127.0.0.1 -c 5 -w /tmp/host_capture.pcap 2>/dev/null &
PID=$!
sleep 1
ping -c 5 127.0.0.1 > /dev/null 2>&1
wait $PID 2>/dev/null
tcpdump -r /tmp/host_capture.pcap 2>/dev/null | wc -l | tr -d " " > /tmp/packet_count.txt

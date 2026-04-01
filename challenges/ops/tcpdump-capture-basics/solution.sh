#!/bin/bash
tcpdump -i any -c 5 -w /tmp/capture.pcap 2>/dev/null &
sleep 2
ping -c 5 127.0.0.1 > /dev/null 2>&1
wait
tcpdump -r /tmp/capture.pcap > /tmp/capture_summary.txt 2>&1

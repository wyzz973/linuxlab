#!/bin/bash
apt-get update -qq && apt-get install -y -qq tcpdump > /dev/null 2>&1
rm -f /tmp/host_capture.pcap /tmp/packet_count.txt

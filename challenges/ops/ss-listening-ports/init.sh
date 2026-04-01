#!/bin/bash
apt-get update -qq && apt-get install -y -qq iproute2 python3 netcat-openbsd > /dev/null 2>&1
rm -f /tmp/listening.txt /tmp/ports.txt
# Start services on multiple ports
python3 -m http.server 8080 &>/dev/null &
python3 -m http.server 9090 &>/dev/null &
nc -l -p 3333 &>/dev/null &
sleep 1

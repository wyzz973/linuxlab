#!/bin/bash
apt-get update -qq && apt-get install -y -qq iproute2 net-tools python3 > /dev/null 2>&1
rm -f /tmp/connections.txt /tmp/conn_stats.txt
python3 -m http.server 8080 &>/dev/null &
sleep 1

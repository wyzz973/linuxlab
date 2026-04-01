#!/bin/bash
apt-get update -qq && apt-get install -y -qq curl python3 > /dev/null 2>&1
rm -f /tmp/custom_ua.txt /tmp/cookies.txt /tmp/verbose.txt
python3 -m http.server 8080 &>/dev/null &
sleep 1

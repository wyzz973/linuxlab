#!/bin/bash
apt-get update -qq && apt-get install -y -qq lsof python3 > /dev/null 2>&1
rm -f /tmp/lsof_listen.txt /tmp/lsof_pid.txt /tmp/lsof_deleted.txt
python3 -m http.server 8080 &>/dev/null &
sleep 1

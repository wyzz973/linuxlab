#!/bin/bash
apt-get update -qq && apt-get install -y -qq at > /dev/null 2>&1
service atd start 2>/dev/null || atd 2>/dev/null
rm -f /tmp/at_help.txt /tmp/at_demo.sh /tmp/at_queue.txt /tmp/at_result.txt

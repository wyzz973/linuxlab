#!/bin/bash
apt-get update -qq && apt-get install -y -qq traceroute iproute2 > /dev/null 2>&1
rm -f /tmp/trace.txt /tmp/gateway.txt

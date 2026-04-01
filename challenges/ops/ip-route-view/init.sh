#!/bin/bash
apt-get update -qq && apt-get install -y -qq iproute2 > /dev/null 2>&1
rm -f /tmp/routes.txt /tmp/route_to_dns.txt /tmp/ip_addrs.txt

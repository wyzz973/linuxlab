#!/bin/bash
apt-get update -qq && apt-get install -y -qq dnsutils host > /dev/null 2>&1
rm -f /tmp/nslookup_fwd.txt /tmp/reverse_dns.txt /tmp/nsswitch_hosts.txt

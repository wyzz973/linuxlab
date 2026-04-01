#!/bin/bash
nslookup localhost > /tmp/nslookup_fwd.txt 2>&1
host 127.0.0.1 > /tmp/reverse_dns.txt 2>&1
grep "^hosts" /etc/nsswitch.conf > /tmp/nsswitch_hosts.txt 2>&1 || grep hosts /etc/nsswitch.conf > /tmp/nsswitch_hosts.txt

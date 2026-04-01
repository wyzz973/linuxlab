#!/bin/bash
ip route show > /tmp/routes.txt 2>&1
ip route get 8.8.8.8 > /tmp/route_to_dns.txt 2>&1
ip addr show > /tmp/ip_addrs.txt 2>&1

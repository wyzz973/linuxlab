#!/bin/bash
apt-get update -qq && apt-get install -y -qq dnsutils > /dev/null 2>&1
rm -f /tmp/dns_servers.txt /tmp/dig_result.txt /tmp/hosts_content.txt

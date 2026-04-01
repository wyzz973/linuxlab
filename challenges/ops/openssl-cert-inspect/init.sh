#!/bin/bash
apt-get update -qq && apt-get install -y -qq openssl > /dev/null 2>&1
rm -f /tmp/server.crt /tmp/server.key /tmp/cert_info.txt /tmp/cert_verify.txt

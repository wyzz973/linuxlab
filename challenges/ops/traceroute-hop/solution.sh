#!/bin/bash
traceroute -m 3 127.0.0.1 > /tmp/trace.txt 2>&1
ip route | grep default | awk "{print \$3}" > /tmp/gateway.txt 2>&1
[ -s /tmp/gateway.txt ] || echo "none" > /tmp/gateway.txt

#!/bin/bash
rm -f /tmp/result.txt
mkdir -p /etc
echo "test config" > /etc/test_challenge.conf 2>/dev/null || true

#!/bin/bash
if grep -q "8443" /home/lab/nginx.conf; then
    if grep -q "listen 80" /home/lab/nginx.conf; then
        echo "FAIL: Port 80 still present"
        exit 1
    fi
    echo "PASS"
    exit 0
else
    echo "FAIL: Port 8443 not found in nginx.conf"
    exit 1
fi

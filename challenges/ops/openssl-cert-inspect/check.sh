#!/bin/bash
if [ ! -f /tmp/server.crt ]; then
    echo "FAIL: /tmp/server.crt not found"
    exit 1
fi
if [ ! -f /tmp/server.key ]; then
    echo "FAIL: /tmp/server.key not found"
    exit 1
fi
if ! openssl x509 -in /tmp/server.crt -noout 2>/dev/null; then
    echo "FAIL: /tmp/server.crt is not a valid certificate"
    exit 1
fi
if [ ! -f /tmp/cert_info.txt ] || [ ! -s /tmp/cert_info.txt ]; then
    echo "FAIL: /tmp/cert_info.txt not found or empty"
    exit 1
fi
if [ ! -f /tmp/cert_verify.txt ] || [ ! -s /tmp/cert_verify.txt ]; then
    echo "FAIL: /tmp/cert_verify.txt not found or empty"
    exit 1
fi
echo "PASS"
exit 0

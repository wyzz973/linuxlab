#!/bin/bash
if [ ! -f /tmp/test_key ]; then
    echo "FAIL: /tmp/test_key not found"
    exit 1
fi
if [ ! -f /tmp/test_key.pub ]; then
    echo "FAIL: /tmp/test_key.pub not found"
    exit 1
fi
if ! grep -q "ssh-ed25519" /tmp/test_key.pub; then
    echo "FAIL: not an Ed25519 key"
    exit 1
fi
if [ ! -f /tmp/pubkey_content.txt ] || [ ! -s /tmp/pubkey_content.txt ]; then
    echo "FAIL: /tmp/pubkey_content.txt not found or empty"
    exit 1
fi
if [ ! -f /tmp/ssh_config.txt ]; then
    echo "FAIL: /tmp/ssh_config.txt not found"
    exit 1
fi
if grep -qi "PermitRootLogin" /tmp/ssh_config.txt; then
    echo "PASS"
    exit 0
fi
echo "FAIL: SSH config missing PermitRootLogin"
exit 1

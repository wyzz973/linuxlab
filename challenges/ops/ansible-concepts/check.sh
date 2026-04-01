#!/bin/bash
if [ ! -f /tmp/ansible_version.txt ]; then
    echo "FAIL: /tmp/ansible_version.txt not found"
    exit 1
fi
if [ ! -f /tmp/inventory.ini ] || [ ! -s /tmp/inventory.ini ]; then
    echo "FAIL: /tmp/inventory.ini not found or empty"
    exit 1
fi
INV=$(cat /tmp/inventory.ini)
if ! echo "$INV" | grep -q "\[webservers\]"; then
    echo "FAIL: missing [webservers] group"
    exit 1
fi
if ! echo "$INV" | grep -q "\[dbservers\]"; then
    echo "FAIL: missing [dbservers] group"
    exit 1
fi
if [ ! -f /tmp/setup.yml ] || [ ! -s /tmp/setup.yml ]; then
    echo "FAIL: /tmp/setup.yml not found or empty"
    exit 1
fi
if grep -qiE "hosts|tasks|nginx" /tmp/setup.yml; then
    echo "PASS"
    exit 0
fi
echo "FAIL: playbook missing required content"
exit 1
